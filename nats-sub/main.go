package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"goProj/dataFactory"
	"log"
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
		//originStr	= fmt.Sprintf("http://localhost:%d", port)
		user 		= "postgres"
		password 	= "postgres"
		dbName		= "testDb"
		dbInfo 		= fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			hostName, port, user, password, dbName)
		//dbURL 		= fmt.Sprintf("postgres://%s:%s@%s:%d/testDb", user, hostName, password, port )
		//natsURL 	= nats.DefaultURL
	)
	db, err := sqlx.Open("postgres", dbInfo)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cr, err := dataFactory.TryGetCreator(msg.Subject, dataFactory.Creators)
	it := cr.Create()
	if err != nil {
		log.Fatal("Unknown item")
	}
	if (msg.Subject != "order") {
		_, err = db.NamedExec(cr.CreateQuery(),it)
	} else {
		_, err = db.Exec(`insert into "Order" (orderuid, entry, internalsignature, payment, items, locale,
                     			 customerid, tracknumber,deliveryservice, shardkey, smid) 
								 values($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11)`,
			it.(dataFactory.Order).OrderUID,
			it.(dataFactory.Order).Entry,
			it.(dataFactory.Order).InternalSignature,
			it.(dataFactory.Order).Payment,
			pq.Array(it.(dataFactory.Order).Items),
			it.(dataFactory.Order).Locale,
			it.(dataFactory.Order).CustomerID,
			it.(dataFactory.Order).TrackNumber,
			it.(dataFactory.Order).DeliveryService,
			it.(dataFactory.Order).Shardkey,
			it.(dataFactory.Order).SmID,)
	}
	if err != nil {
		fmt.Println(err)
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