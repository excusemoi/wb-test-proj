package dataFactory

import (
	"database/sql/driver"
	"fmt"
)

type OutputOrder struct {
	OrderUID string `json:"order_uid"`
	Entry string `json:"entry"`
	TotalPrice int `json:"total_price"`
	CustomerID string `json:"customer_id"`
	TrackNumber string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
}

type IData interface {
}

type Order struct {
	OrderUID string `json:"order_uid" db:"orderuid"`
	Entry string `json:"entry" db:"entry"`
	InternalSignature string `json:"internal_signature" db:"internalsignature"`
	Payment Payment `json:"payment" db:"payment"`
	Items []Item `json:"dataFactory" db:"dataFactory"`
	Locale string `json:"locale" db:"locale"`
	CustomerID string `json:"customer_id" db:"customerid"`
	TrackNumber string `json:"track_number" db:"tracknumber"`
	DeliveryService string `json:"delivery_service" db:"deliveryservice"`
	Shardkey string `json:"shardkey" db:"shardkey"`
	SmID int `json:"sm_id" db:"smid"`
	PaymentID int `json:"Payment_id" db:"PaymentID"`
	IData
}



type Payment struct {
	IData
	driver.Valuer
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
func (p Payment) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s, %s, %s, %d, %d, %s, %d, %d)",
		p.Transaction,
		p.Currency,
		p.Provider,
		p.Amount,
		p.PaymentDt,
		p.Bank,
		p.DeliveryCost,
		p.GoodsTotal), nil
}

type Item struct {
	IData
	driver.Valuer
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

func (i Item) Value() (driver.Value, error) {
	return fmt.Sprintf("(%d, %d, %s, %s, %d, %s, %d, %d, %s)",
		i.ChrtID,
		i.Price,
		i.Rid,
		i.Name,
		i.Sale,
		i.Size,
		i.TotalPrice,
		i.NmID,
		i.Brand), nil
}

