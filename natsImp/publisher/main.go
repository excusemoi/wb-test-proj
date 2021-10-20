package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"goProj/config"
	"goProj/dataFactory"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg := config.Get("../../config/config.json")

	clusterUrls := strings.Join(cfg.ClusterUrls, ", ")
	sc, err := stan.Connect(cfg.ClusterId,
		"id1",
		stan.NatsURL(clusterUrls),
		stan.NatsOptions(
			nats.ReconnectWait(time.Second*4),
			nats.Timeout(time.Second*4)),
	)

	if err != nil {
		return err
	}

	ordersCreator := dataFactory.OrderCreator{}
	order := ordersCreator.Create(rand.Intn(10000))

	message, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		return err
	}

	err = sc.Publish(cfg.Subject, message)

	log.Printf("Published random %s\n", cfg.Subject)

	return nil
}