package api

import (
	"database/sql"
	"fmt"
	"net/http"
	apimodels "promotions/api/models"
	persistance "promotions/persistance"
	db_entities "promotions/persistance/entities"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPromotionByPromotionID(c *gin.Context, db *sql.DB) {
	promotion_id := c.Param("promotion_id")
	promotion, dbError := persistance.GetPromotionByPromotionID(promotion_id, db)
	if dbError != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": dbError.Error()})
		return
	}

	getPromotionResponse := parseDbEntityToGetPromotionResponse(promotion)

	c.JSON(http.StatusOK, getPromotionResponse)
}

func CreatePromotion(c *gin.Context, db *sql.DB) {
	var promotion apimodels.CreatePromotionRequest
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbPromotion, parseErr := parseToDbEntity(&promotion)
	var message string
	if parseErr != nil {
		message = "Error during parsing promotion to db entity."
		c.JSON(http.StatusBadRequest, generateCreatePromotionResponse(0, message, promotion))
		return
	}

	id, dbErr := persistance.AddPromotion(dbPromotion, db)
	if dbErr != nil {
		message = "Error during saving promotion to db."
		c.JSON(http.StatusBadRequest, generateCreatePromotionResponse(0, message, promotion))
		return
	}

	message = "Promotion Created Successfully!"
	createPromotionResponse := generateCreatePromotionResponse(id, message, promotion)

	c.JSON(http.StatusCreated, createPromotionResponse)
}

// Utility functions

func parseToDbEntity(apiPromotion *apimodels.CreatePromotionRequest) (db_entities.Promotion, error) {
	expDate, parseErr := parseStringToTime(apiPromotion.ExpirationDate, "2006-01-02 15:04:05 -0700 MST")
	if parseErr != nil {
		return db_entities.Promotion{}, parseErr
	}

	return db_entities.Promotion{
		PromotionId:    apiPromotion.PromotionId,
		Price:          apiPromotion.Price,
		ExpirationDate: expDate,
	}, nil
}

func parseStringToTime(dateString string, layout string) (time.Time, error) {
	date, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Failed to parse date string:", err)
		return time.Time{}, err
	}
	return date, nil
}

func parseTimeToString(time time.Time, layout string) string {
	return time.Format(layout)
}

func parseDbEntityToGetPromotionResponse(dbEntity db_entities.Promotion) apimodels.GetPromotionResponse {
	return apimodels.GetPromotionResponse{
		ID:             dbEntity.ID,
		PromotionId:    dbEntity.PromotionId,
		Price:          dbEntity.Price,
		ExpirationDate: parseTimeToString(dbEntity.ExpirationDate, "2006-01-02 15:04:05 -0700 MST"),
	}
}

func generateCreatePromotionResponse(id int64, message string,
	createPromrequest apimodels.CreatePromotionRequest) apimodels.CreatePromotionResponse {
	return apimodels.CreatePromotionResponse{
		Message: message,
		PromotionResponse: apimodels.GetPromotionResponse{
			ID:             id,
			PromotionId:    createPromrequest.PromotionId,
			Price:          createPromrequest.Price,
			ExpirationDate: createPromrequest.ExpirationDate,
		},
	}
}
