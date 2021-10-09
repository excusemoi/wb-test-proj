package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"runtime"
	"time"
	"../items"
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

	nc.Subscribe(subj, func(msg *nats.Msg) {
		fmt.Printf("Item: %s Uid: %s\n", msg.Subject, msg.Data)

	})
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
		originStr	= fmt.Sprintf("http://localhost:%d", port)
		user 		= "postgres"
		password 	= "postgres"
		dbName		= "testDb"
		dbInfo 		= fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			hostName, port, user, password, dbName)
		dbURL 		= fmt.Sprintf("postgres://%s:%s@%s:%d/testDb", user, hostName, password, port )
		natsURL 	= nats.DefaultURL
	)
	fmt.Printf("Data:\n\t%s\n\t%s\n\t%s\n\t%s", originStr, dbURL, natsURL, dbInfo)
	fmt.Println("Trying to connect to db...")
	db, err := sql.Open("postgres", dbInfo)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var typeOfItem, uId = msg.Subject, msg.Data
	fmt.Printf("Trying to add %s %s\n", typeOfItem, uId)
	cr, err := items.TryGetCreator(typeOfItem, items.Creators)
	if err != nil {
		log.Fatal("Unknown item")
	}
	db.Exec("insert into $1 values($2)", typeOfItem, cr.Create())
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