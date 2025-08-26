package models

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid"          gorm:"primaryKey"`
	TrackNumber       string    `json:"track_number"       gorm:"not null;unique"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"           gorm:"foreignKey:OrderUID;constraint:OnDelete:CASCADE;" musttag:"true"`
	Payment           Payment   `json:"payment"            gorm:"foreignKey:OrderUID;constraint:OnDelete:CASCADE;" musttag:"true"`
	Items             []Item    `json:"items"              gorm:"foreignKey:OrderUID;constraint:OnDelete:CASCADE;" musttag:"true"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	OrderUID string `gorm:"primaryKey" json:"-"`
	Name     string `                  json:"name"`
	Phone    string `                  json:"phone"`
	Zip      string `                  json:"zip"`
	City     string `                  json:"city"`
	Address  string `                  json:"address"`
	Region   string `                  json:"region"`
	Email    string `                  json:"email"`
}

type Payment struct {
	OrderUID     string  `gorm:"primaryKey"      json:"-"`
	Transaction  string  `gorm:"unique;not null" json:"transaction"`
	RequestID    string  `                       json:"request_id"`
	Currency     string  `                       json:"currency"`
	Provider     string  `                       json:"provider"`
	Amount       float64 `                       json:"amount"`
	PaymentDt    int64   `                       json:"payment_dt"`
	Bank         string  `                       json:"bank"`
	DeliveryCost float64 `                       json:"delivery_cost"`
	GoodsTotal   float64 `                       json:"goods_total"`
	CustomFee    float64 `                       json:"custom_fee"`
}

type Item struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"-"`
	OrderUID    string  `gorm:"index;not null"           json:"-"`
	ChrtID      int     `gorm:"not null"                 json:"chrt_id"`
	TrackNumber string  `                                json:"track_number"`
	Price       float64 `                                json:"price"`
	Rid         string  `                                json:"rid"`
	Name        string  `                                json:"name"`
	Sale        int     `                                json:"sale"`
	Size        string  `                                json:"size"`
	TotalPrice  float64 `                                json:"total_price"`
	NmID        int     `                                json:"nm_id"`
	Brand       string  `                                json:"brand"`
	Status      int     `                                json:"status"`
}
