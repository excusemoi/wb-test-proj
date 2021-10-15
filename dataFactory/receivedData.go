package dataFactory

type IData interface {
}

type OutputOrder struct {
	OrderUID string `json:"order_uid"`
	Entry string `json:"entry"`
	TotalPrice int `json:"total_price"`
	CustomerID string `json:"customer_id"`
	TrackNumber string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
}

type Order struct {
	IData
	OrderUID string `json:"order_uid" db:"orderuid"`
	Entry string `json:"entry" db:"entry"`
	InternalSignature string `json:"internal_signature" db:"internalsignature"`
	Payment Payment `json:"payment" db:"payment"`
	Items []Item `json:"items" db:"items"`
	Locale string `json:"locale" db:"locale"`
	CustomerID string `json:"customer_id" db:"customerid"`
	TrackNumber string `json:"track_number" db:"tracknumber"`
	DeliveryService string `json:"delivery_service" db:"deliveryservice"`
	Shardkey string `json:"shardkey" db:"shardkey"`
	SmID int `json:"sm_id" db:"smid"`
	PaymentID int `json:"Payment_id" db:"paymentid"`
}

type Payment struct {
	IData
	Transaction string `json:"transaction" db:"transaction"`
	Currency string `json:"currency" db:"currency"`
	Provider string `json:"provider" db:"provider"`
	Amount int `json:"amount" db:"amount"`
	PaymentDt int `json:"payment_dt" db:"paymentdt"`
	Bank string `json:"bank" db:"bank"`
	DeliveryCost int `json:"delivery_cost" db:"deliverycost"`
	GoodsTotal int `json:"goods_total" db:"goodstotal"`
	PaymentID int `db:"PaymentID"`
}

type Item struct {
	IData
	ChrtID string `json:"chrt_id" db:"chrtid"`
	Price int `json:"price" db:"price"`
	Rid string `json:"rid" db:"rid"`
	Name string `json:"name" db:"name"`
	Sale int `json:"sale" db:"sale"`
	Size string `json:"size" db:"size"`
	TotalPrice int `json:"total_price" db:"totalprice"`
	NmID int `json:"nm_id" db:"nmid"`
	Brand string `json:"brand" db:"brand"`
	ItemID int `db:"ItemID"`
}