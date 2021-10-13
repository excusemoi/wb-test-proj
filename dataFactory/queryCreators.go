package dataFactory

var ItemQuery = `insert into "Items" (chrtid, price, rid, name, sale, size, totalprice, nmid, brand, ItemID) 
			values(:chrtid, :price, :rid, :name, :sale, :size, :totalprice, :nmid, :brand, :ItemID)`
var PaymentQuery = `insert into "Payments" (transaction, currency, provider, amount, paymentdt, bank, deliverycost, goodstotal, PaymentID) 
			values(:transaction, :currency, :provider, :amount, :paymentdt, :bank, :deliverycost, :goodstotal, :PaymentID)`
