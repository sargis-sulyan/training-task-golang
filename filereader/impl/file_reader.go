package filereader

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	file_models "promotions/filereader/models"
	persistance "promotions/persistance"
	db_entities "promotions/persistance/entities"
	"strconv"
	"sync"
	"time"

	"database/sql"
)

// ReadAndSaveFileData reads the data from csv file and saves it concurrently
func ReadAndSaveFileData(db *sql.DB) {
	log.Println("Truncating table data ...")
	truncateError := truncatePromotionsTable(db)

	if truncateError != nil {
		log.Println(truncateError.Error())
		return
	}
	log.Println("Truncated table data!")

	csvFile, openErr := os.Open("../resources/promotions.csv")
	if openErr != nil {
		log.Println(openErr.Error())
		return
	}
	defer csvFile.Close()

	timeBeforeCall := time.Now()
	concurrentReadAndSaveToDbC(csvFile, db)
	timeDiff := time.Now().Sub(timeBeforeCall)

	log.Println("End concurrent reading of csv file and saving data to database, time : ", timeDiff)
}

// read and save data with Worker pools
func concurrentReadAndSaveToDbC(csvFile *os.File, db *sql.DB) {
	log.Println("Start reading data from file ...")
	fcsv := csv.NewReader(csvFile)
	numWps := 100
	jobs := make(chan []string, numWps)
	res := make(chan *file_models.Promotion)

	var wg sync.WaitGroup
	worker := func(jobs <-chan []string, results chan<- *file_models.Promotion) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				results <- parseStruct(job)
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs, res)
		}()
	}

	go func() {
		for {
			rStr, err := fcsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("ERROR: ", err.Error())
				break
			}
			jobs <- rStr
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	log.Println("Start saving data to db ...")
	for r := range res {
		saveToDB(r, db)
	}
	log.Println("Finished saving data to db!")
	persistance.CloseDB(db)
}

func truncatePromotionsTable(db *sql.DB) error {
	return persistance.TruncatePromotions(db)

}

func saveToDB(fileEntity *file_models.Promotion, db *sql.DB) {
	dbEntity, parseErr := parseToDbEntity(fileEntity)
	if parseErr != nil {
		log.Println(parseErr)
		return
	}

	_, savErr := persistance.AddPromotion(dbEntity, db)
	if savErr != nil {
		log.Println(savErr)
		return
	}
}

func parseStruct(data []string) *file_models.Promotion {
	price, _ := strconv.ParseFloat(data[1], 64)
	return &file_models.Promotion{
		Id:             data[0],
		Price:          price,
		ExpirationDate: data[2],
	}
}

func parseToDbEntity(fileEntity *file_models.Promotion) (db_entities.Promotion, error) {
	expirationDate, parseErr := parseStringToTime(fileEntity.ExpirationDate, "2006-01-02 15:04:05 -0700 MST")
	if parseErr != nil {
		return db_entities.Promotion{}, parseErr
	}
	return db_entities.Promotion{
		PromotionId:    fileEntity.Id,
		Price:          fileEntity.Price,
		ExpirationDate: expirationDate,
	}, nil
}

func parseStringToTime(dateString string, layout string) (time.Time, error) {
	date, err := time.Parse(layout, dateString)
	if err != nil {
		log.Println("Failed to parse date string:", err)
		return time.Time{}, err
	}
	return date, nil
}
