package persistance

import (
	"time"
)

type Promotion struct {
	ID             int64
	PromotionId    string
	Price          float64
	ExpirationDate time.Time
}
