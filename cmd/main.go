package main

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"goProj/dataFactory"
	"log"
	"net/http"
	//"github.com/blockloop/scan"
	//"os"
)

var port = "3000"

func serveFiles(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		path := request.URL.Path
		fmt.Println(path)
		if path == "/" {
			path = "../web/static/main.html"
		} else {
			path = "." + path
		}
		http.ServeFile(writer, request, path)
	case "POST":
		id := request.FormValue("message")
		var (
			dbPort 		= 5432
			hostName 	= "localhost"
			user 		= "postgres"
			password 	= "postgres"
			dbName		= "testDb"
			dbInfo 		= fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s dbname=%s sslmode=disable",
				hostName, dbPort, user, password, dbName)
		)
		db, _ := sqlx.Open("postgres", dbInfo)
		err := db.Ping()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		res, err := db.Query(fmt.Sprintf(`select * from "Order" where OrderUID = '%s' limit 1;`, id))
		if err != nil {
			log.Fatal(err)
		}
		o := dataFactory.Order{}
		var deliveryCost int
		var totalPrice int

		for res.Next() {
			if err = res.Scan(&o.OrderUID, &o.Entry, &o.InternalSignature,
				&o.Locale, &o.CustomerID, &o.TrackNumber, &o.DeliveryService,
				&o.Shardkey, &o.SmID, &o.PaymentID); err != nil {
				fmt.Println(err.Error())
				return
			}

		}

		res, err = db.Query(fmt.Sprintf(`select p.paymentid from "Payments" as p where paymentid = '%d';`, o.PaymentID))
		if err != nil {
			log.Fatal(err)
		}
		for res.Next() {
			res.Scan(&deliveryCost)
		}

		res, err = db.Query(fmt.Sprintf(`select i.chrtid, sum(i.totalprice) from "Items" as i where chrtid = '%s'
													group by i.chrtid;`, o.OrderUID))
		if err != nil {
			log.Fatal(err)
		}
		for res.Next() {
			res.Scan(&o.OrderUID, &totalPrice)
		}
		totalPrice += deliveryCost
		outputOrderBytes, err := json.MarshalIndent(dataFactory.OutputOrder{
			OrderUID:        o.OrderUID,
			Entry:           o.Entry,
			TotalPrice:      totalPrice,
			CustomerID:      o.CustomerID,
			TrackNumber:     o.TrackNumber,
			DeliveryService: o.DeliveryService,
		}, "", "\t")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(writer, string(outputOrderBytes))

	default:
		fmt.Fprintf(writer,"Request type other than GET or Post not supported")
	}
}

func main() {
	fmt.Println("Listening on port :" + port)
	http.HandleFunc("/", serveFiles)
	http.ListenAndServe(":" + port, nil)
}

