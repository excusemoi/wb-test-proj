package dataFactory

func (ic ItemCreator) CreateQuery() string {
	return `insert into "Items" (chrtid, price, rid, name, sale, size, totalprice, nmid, brand) 
			values(:chrtid, :price, :rid, :name, :sale, :size, :totalprice, :nmid, :brand)`
}

func (pc PaymentCreator) CreateQuery() string {
	return `insert into "Payments" (transaction, currency, provider, amount, paymentdt, bank, deliverycost, goodstotal) 
			values(:transaction, :currency, :provider, :amount, :paymentdt, :bank, :deliverycost, :goodstotal)`
}
