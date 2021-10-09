package items

import (
	"../resources"
	"errors"
	"math/rand"
	"time"
)

var Creators = map[string]ICreator{
	"item" : ItemCreator{},
	"payment" : PaymentCreator{},
}

func TryGetCreator(productName string, creators map[string]ICreator) (ICreator, error) {
	if cr, exist := Creators[productName]; exist {
		return cr, nil
	}
	return nil, errors.New("Creator doesn't exist")
}

type ICreator interface {
	Create() (data.IData)
}

type OrderCreator struct {
}

type ItemCreator struct {
}

type PaymentCreator struct {
}

func (ic ItemCreator) Create() data.IData {
	rand.Seed(time.Now().UnixNano())
	var price = randInt(minPrice, maxPrice + 1)
	var sale = rand.Intn(minPrice)
	var totalPrice = price - sale
	return data.Item{
		ChrtID:     rand.Intn(100),
		Price:      price,
		Rid:        randomString(5),
		Name:       items[rand.Intn(len(items))],
		Sale:       sale,
		Size:       randomString(5),
		TotalPrice: totalPrice,
		NmID:       rand.Intn(5),
		Brand:      brands[rand.Intn(len(brands))],
	}
}

func (pc PaymentCreator) Create() data.IData {
	rand.Seed(time.Now().UnixNano())
	return data.Payment{
		Transaction:  transactions[rand.Intn(len(transactions))],
		Currency:     currencys[rand.Intn(len(currencys))],
		Provider:     providers[rand.Intn(len(providers))],
		Amount:       rand.Intn(5),
		PaymentDt:    rand.Intn(5),
		Bank:         banks[rand.Intn(len(banks))],
		DeliveryCost: randInt(minPrice, maxPrice),
		GoodsTotal:   randInt(0,10),
	}
}

//func (oc OrderCreator) create() data.IData {
	//rand.Seed(time.Now().UnixNano())
	//var	o = data.Order{}
//}


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

var currencys = [...]string {
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