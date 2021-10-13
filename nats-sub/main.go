package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"goProj/dataFactory"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var showTime = flag.Bool("t", false, "Display timestamps")

	log.SetFlags(0)
	flag.Parse()

	args := flag.Args()

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)

	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj := args[0]

	nc.Subscribe(subj, handleMessage)

	nc.Flush()

	if err = nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}

func handleMessage(msg *nats.Msg){
	var (
		port 		= 5432
		hostName 	= "localhost"
		user 		= "postgres"
		password 	= "postgres"
		dbName		= "testDb"
		dbInfo 		= fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			hostName, port, user, password, dbName)
	)

	db, err := sqlx.Open("postgres", dbInfo)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	id := rand.Intn(10000)
	cr, err := dataFactory.TryGetCreator(msg.Subject, dataFactory.Creators)
	it := cr.Create(id)
	if err != nil {
		log.Fatal("Unknown item")
	}
	switch msg.Subject {
	case "order":
		db.NamedExec(dataFactory.PaymentQuery, it.(dataFactory.Order).Payment)
		_, err = db.Exec(`insert into "Order" (orderuid, entry, internalsignature, locale,
                     			 customerid, tracknumber,deliveryservice, shardkey, smid, paymentid) 
								 values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			it.(dataFactory.Order).OrderUID,
			it.(dataFactory.Order).Entry,
			it.(dataFactory.Order).InternalSignature,
			it.(dataFactory.Order).Locale,
			it.(dataFactory.Order).CustomerID,
			it.(dataFactory.Order).TrackNumber,
			it.(dataFactory.Order).DeliveryService,
			it.(dataFactory.Order).Shardkey,
			it.(dataFactory.Order).SmID,
			it.(dataFactory.Order).PaymentID)

		for _, item := range it.(dataFactory.Order).Items {
			item.ChrtID = it.(dataFactory.Order).OrderUID
			_, err = db.NamedExec(dataFactory.ItemQuery, item)
		}
	case "payment":
		_,err = db.NamedExec(dataFactory.PaymentQuery, it)
	default:
		fmt.Println("Can't add object to db")
	}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}