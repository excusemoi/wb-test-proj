package server

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"goProj/config"
	"goProj/dataFactory"
	"goProj/db"
	"log"
	"net/http"
)

type Server struct {
	cfg 			*config.Config
	db 				*sqlx.DB
	serverPort 		string
}

func InitServer() (*Server, error){
	cfg := config.Get("../config/config.json")
	pgDb, err := db.Dial(cfg)
	if err != nil {
		return nil, err
	}
	return &Server{
		cfg:        cfg,
		db:         pgDb,
		serverPort: "3000",
	}, nil
}

func (s *Server) Run(){
	fmt.Println("Listening on port :" + s.serverPort)
	http.HandleFunc("/", s.HandleFunction)
	http.ListenAndServe(":" + s.serverPort, nil)
}

func (s *Server) HandleFunction(writer http.ResponseWriter, request *http.Request){
	switch request.Method {
	case "GET":
		s.get(writer, request)
	case "POST":
		s.post(writer, request)
	default:
		log.Fatal("Request type other than GET or Post not supported\n")
	}
}

func (*Server) get(writer http.ResponseWriter, request *http.Request){
	path := request.URL.Path
	if path == "/" {
		path = "./static/main.html"
	} else {
		path = "." + path
	}
	http.ServeFile(writer, request, path)
}

func (s *Server) post(writer http.ResponseWriter, request *http.Request) {
	id := request.FormValue("message")
	res, err := s.db.Query(fmt.Sprintf(`select * from "Order" where OrderUID = '%s' limit 1;`, id))
	if err != nil {
		log.Fatal(err)
	}
	o := dataFactory.Order{}
	var deliveryCost int
	var totalPrice int

	for res.Next() {
		if err = res.Scan(&o.OrderUID, &o.Entry, &o.InternalSignature,
			&o.Locale, &o.CustomerID, &o.TrackNumber, &o.DeliveryService,
			&o.Shardkey, &o.SmID, &o.PaymentID);
		err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	res, err = s.db.Query(fmt.Sprintf(`select p.paymentid from "Payments" as p where paymentid = '%d';`,
																										o.PaymentID))
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		res.Scan(&deliveryCost)
	}

	res, err = s.db.Query(fmt.Sprintf(`select i.chrtid, sum(i.totalprice) from "Items" as i where chrtid = '%s'
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
	if o.OrderUID != "" {
		fmt.Fprintf(writer, string(outputOrderBytes))
	}
}

