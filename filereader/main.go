package main

import (
	"fmt"
	"log"
	filereader "promotions/filereader/impl"
	persistance "promotions/persistance"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	//create the scheduler
	loc, locErr := time.LoadLocation("EST")
	if locErr != nil {
		panic(locErr)
	}

	scheduler := gocron.NewScheduler(loc)
	scheduler.Every(30).Minute().Do(readAndSaveFileData)
	scheduler.StartBlocking()
}

func readAndSaveFileData() {
	// get db connection
	log.Println("Getting db connection")
	db, err := persistance.GetMySqlConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// read and save file data
	filereader.ReadAndSaveFileData(db)
}
