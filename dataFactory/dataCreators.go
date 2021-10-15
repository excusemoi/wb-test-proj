package dataFactory

import (
	"math/rand"
	"time"
)

type ICreator interface {
	Create(id int) IData
}

type OrderCreator struct {
}

type ItemCreator struct {
}

type PaymentCreator struct {
}

func (ic ItemCreator) Create(id int) IData {
	rand.Seed(time.Now().UnixNano())
	var price = randInt(minPrice, maxPrice + 1)
	var sale = rand.Intn(minPrice)
	var totalPrice = price - sale
	return Item{
		ChrtID:     randomString(5),
		Price:      price,
		Rid:        randomString(5),
		Name:       items[rand.Intn(len(items))],
		Sale:       sale,
		Size:       randomString(5),
		TotalPrice: totalPrice,
		NmID:       rand.Intn(5),
		Brand:      brands[rand.Intn(len(brands))],
		ItemID: id,
	}
}

func (pc PaymentCreator) Create(id int) IData {
	rand.Seed(time.Now().UnixNano())
	return Payment{
		Transaction:  transactions[rand.Intn(len(transactions))],
		Currency:     currencies[rand.Intn(len(currencies))],
		Provider:     providers[rand.Intn(len(providers))],
		Amount:       rand.Intn(5),
		PaymentDt:    rand.Intn(5),
		Bank:         banks[rand.Intn(len(banks))],
		DeliveryCost: randInt(minPrice, maxPrice),
		GoodsTotal:   randInt(0,10),
		PaymentID:    id,
	}
}

func (oc OrderCreator) Create(id int) IData {
	rand.Seed(time.Now().UnixNano())
	pc := PaymentCreator{}
	ic := ItemCreator{}
	var o = Order{
		IData:             nil,
		OrderUID:          randomString(5),
		Entry:             randomString(5),
		InternalSignature: randomString(5),
		Payment:           pc.Create(id).(Payment),
		Items:             []Item{},
		Locale:            randomString(5),
		CustomerID:        randomString(5),
		TrackNumber:       randomString(5),
		DeliveryService:   deliveryServices[rand.Intn(len(deliveryServices))],
		Shardkey:          randomString(5),
		SmID:              rand.Intn(5),
		PaymentID: id,
	}
	itemsAmount := randInt(1, 5)
	for i := 0; i < itemsAmount; i++ {
		o.Items = append(o.Items, ic.Create(id).(Item))
		o.Items[i].ChrtID = o.OrderUID
	}
	o.Payment.Amount = itemsAmount
	return o
}

const minPrice = 20
const maxPrice = 10000

var banks = [...]string{
	"Allahabad Bank",
	"Andhra Bank",
	"Axis Bank",
	"Bank of Bahrain and Kuwait",
	"Bank of Baroda - Corporate Banking",
	"Bank of Baroda - Retail Banking",
	"Bank of India",
	"Bank of Maharashtra",
	"Canara Bank",
	"Central Bank of India",
	"City Union Bank",
}

var providers = [...]string{
	"PayPal",
	"Due",
	"Stripe",
	"Flagship Merchant Services",
	"Payline Data",
	"Square",
	"Adyen",
	"BirPay",
}

var currencies = [...]string {
	"RUB",
	"EUR",
	"AUD",
	"BRL",
	"BGN",
	"KHR",
	"CVE",
	"KYD",
	"XAF",
	"CLP",
}

var transactions = [...]string {
	"Cash",
	"Personal Cheque",
	"Debit Card",
	"Credit Card",
}

var items = [...]string {
	"Water",
	"Book",
	"Shoes",
	"Ladder",
	"Paper",
	"Pen",
	"Shirt",
	"Cake",
	"Pencil",
	"PC",
}

var brands = [...]string {
	"Adidas",
	"Coca-cola",
	"Erich Krause",
	"Subway",
	"Chanel",
	"Nile",
	"Samsung",
	"MacDonalds",
	"KFC",
}

var deliveryServices = [...]string {
	"Meituan",
	"Uber Eats",
	"Delivery Hero",
	"DoorDash",
	"Grubhub",
	"Deliveroo",
	"Just Eat",
	"Postmates",
	"Swiggy",
	"Zomato",
}
