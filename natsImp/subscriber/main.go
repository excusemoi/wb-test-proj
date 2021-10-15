package main

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"goProj/config"
	"goProj/dataFactory"
	"goProj/db"
	"goProj/natsImp"
	"log"
	"runtime"
)

type subscriber struct {
	ns *natsImp.Nats
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	runtime.Goexit()
}

func run() error {
	ni, err := natsImp.InitNats()
	if err != nil {
		return err
	}
	sub := subscriber{ns: ni}
	sub.ns.Conn, err = nats.Connect(nats.DefaultURL, ni.Options...)

	cfg := config.Get("../../config/config.json")
	dbPg, err := db.Dial(cfg)
	if err != nil {
		return err
	}

	_, err = sub.ns.Conn.Subscribe("order", func(msg *nats.Msg) {
		o := dataFactory.Order{}
		err = json.Unmarshal(msg.Data, &o)
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = dbPg.NamedExec(dataFactory.PaymentQuery, o.Payment)
		_, err = dbPg.NamedExec(dataFactory.OrderQuery, o)
		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, item := range o.Items {
			_, err = dbPg.NamedExec(dataFactory.ItemQuery, item)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	})

	err = sub.ns.Conn.Flush()

	if err = sub.ns.Conn.LastError(); err != nil {
		return err
	}

	log.Printf("Listening order")

	return nil
}

