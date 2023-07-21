package filereader

type Promotion struct {
	Id             string  `json:"id"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}
