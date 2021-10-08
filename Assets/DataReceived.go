package main

type Order struct {
	OrderUID string `json:"order_uid"`
	Entry string `json:"entry"`
	InternalSignature string `json:"internal_signature"`
	Payment Payment `json:"payment"`
	Items []Items `json:"items"`
	Locale string `json:"locale"`
	CustomerID string `json:"customer_id"`
	TrackNumber string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
	Shardkey string `json:"shardkey"`
	SmID int `json:"sm_id"`
}
type Payment struct {
	Transaction string `json:"transaction"`
	Currency string `json:"currency"`
	Provider string `json:"provider"`
	Amount int `json:"amount"`
	PaymentDt int `json:"payment_dt"`
	Bank string `json:"bank"`
	DeliveryCost int `json:"delivery_cost"`
	GoodsTotal int `json:"goods_total"`
}
type Items struct {
	ChrtID int `json:"chrt_id"`
	Price int `json:"price"`
	Rid string `json:"rid"`
	Name string `json:"name"`
	Sale int `json:"sale"`
	Size string `json:"size"`
	TotalPrice int `json:"total_price"`
	NmID int `json:"nm_id"`
	Brand string `json:"brand"`
}


