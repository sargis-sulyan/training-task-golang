package persistance

import (
	"database/sql"
	"fmt"
	"log"
	db_entities "promotions/persistance/entities"
)

// 1.GET

// GetPromotionByID queries for the promotion with the specified ID.
func GetPromotionByID(id int, db *sql.DB) (db_entities.Promotion, error) {
	// An promotion to hold data from the returned row.
	var pro db_entities.Promotion

	row := db.QueryRow("SELECT * FROM promotions WHERE id = ?", id)
	if err := row.Scan(&pro.ID, &pro.PromotionId, &pro.Price, &pro.ExpirationDate); err != nil {
		if err == sql.ErrNoRows {
			return pro, fmt.Errorf("promotionsById %d: no such promotion", id)
		}
		return pro, fmt.Errorf("promotionsById %d: %v", id, err)
	}
	return pro, nil
}

// GetPromotionByPromotionID queries for the promotion with the specified promotion_id.
func GetPromotionByPromotionID(promotion_id string, db *sql.DB) (db_entities.Promotion, error) {
	// An promotion to hold data from the returned row.
	var pro db_entities.Promotion

	row := db.QueryRow("SELECT * FROM promotions WHERE promotion_id = ?", promotion_id)
	if err := row.Scan(&pro.ID, &pro.PromotionId, &pro.Price, &pro.ExpirationDate); err != nil {
		if err == sql.ErrNoRows {
			return pro, fmt.Errorf("promotionsByPromotionId %d: no such promotion", promotion_id)
		}
		return pro, fmt.Errorf("promotionsByPromotionId %d: %v", promotion_id, err)
	}
	return pro, nil
}

// 2.CREATE

// AddPromotion adds the specified promotion to the database,
// returning the promotion ID of the new entry
func AddPromotion(pro db_entities.Promotion, db *sql.DB) (int64, error) {
	result, err := db.Exec("INSERT INTO promotions (promotion_id, price, expiration_date) VALUES (?, ?, ?)", pro.PromotionId, pro.Price, pro.ExpirationDate)
	if err != nil {
		return 0, fmt.Errorf("addPromotion: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addPromotion: %v", err)
	}
	return id, nil
}

// AddPromotions adds the specified promotions to the database
func AddPromotions(proms []db_entities.Promotion, db *sql.DB) (int64, error) {
	for p := range proms {
		albID, err := AddPromotion(proms[p], db)
		if err != nil {
			log.Fatal(err)
			return 0, err
		}
		fmt.Printf("ID of added promotion: %v\n", albID)
	}
	return 0, nil
}

// 3.DELETE

// DeletePromotionByID deletes promotion by specified id
func DeletePromotionByID(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM promotions WHERE id = ?", id)
	return err
}

// DeletePromotionByID deletes promotion by specified promotion_id.
func DeletePromotionByPromotionId(db *sql.DB, promotion_id string) error {
	_, err := db.Exec("DELETE FROM promotions WHERE promotion_ = ?", promotion_id)
	return err
}

// 4.TRUNCATE
// TruncatePromotions truncates the promotions table
func TruncatePromotions(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE promotions;")
	return err
}
