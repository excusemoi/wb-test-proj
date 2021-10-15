package natsImp

import (
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Nats struct {
	Conn *nats.Conn
	Options []nats.Option
}

func InitNats() (*Nats, error) {
	var (
		userCreds = flag.String("creds", "", "User Credentials File")
		nkeyFile = flag.String("nkey", "", "NKey Seed File")
		tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
		tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
		tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	)
	log.SetFlags(0)
	flag.Parse()
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
	return &Nats{Conn: &nats.Conn{}, Options: opts}, nil
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