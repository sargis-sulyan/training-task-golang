package main

import (
	"database/sql"
	"fmt"
	controllers "promotions/api/controllers"
	persistance "promotions/persistance"

	"github.com/gin-gonic/gin"
)

func main() {
	err := persistance.InitMySqlConnectionPool(1000)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to DB!")

	db := persistance.GetConnection()
	defer db.Close()

	router := setupRouter(db)
	router.Run(":9090")

}

func setupRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()

	router.GET("/promotions/:promotion_id", func(c *gin.Context) {
		controllers.GetPromotionByPromotionID(c, db)
	})
	router.POST("/promotions", func(c *gin.Context) {
		controllers.CreatePromotion(c, db)
	})

	//router.GET("/promotions", controllers.GetTodos)
	// router.PUT("/promotions/:id", controllers.UpdateTodo)
	// router.DELETE("/promotions/:id", controllers.DeleteTodo)

	return router
}
