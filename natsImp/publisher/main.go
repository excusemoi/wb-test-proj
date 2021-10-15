package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"goProj/dataFactory"
	"goProj/natsImp"
	"log"
	"math/rand"
)

type publisher struct {
	ns *natsImp.Nats
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ni, err := natsImp.InitNats()
	if err != nil {
		return err
	}
	pub := publisher{ns: ni}
	pub.ns.Conn, err = nats.Connect(nats.DefaultURL, pub.ns.Options...)
	defer pub.ns.Conn.Close()

	ordersCreator := dataFactory.OrderCreator{}
	order := ordersCreator.Create(rand.Intn(10000))

	message, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		return err
	}

	err = pub.ns.Conn.Publish("order", message)
	err = pub.ns.Conn.Flush()

	if err = pub.ns.Conn.LastError(); err != nil {
		return err
	}  else {
		log.Printf("Published random order\n")
	}
	return nil
}
