package subscriber

import (
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"goProj/config"
	"goProj/dataFactory"
	"goProj/db"
	"log"
	"os"
	"strings"
	"time"
)

func Run(pathToCfg string) error {
	if len(os.Args) < 2 {
		return errors.New("Please indicate the path to the config file")
	}

	cfg := config.Get(pathToCfg)
	dbPg, err := db.Dial(cfg)
	if err != nil {
		return err
	}

	clusterUrls := strings.Join(cfg.ClusterUrls, ", ")
	sc, err := stan.Connect(cfg.ClusterId,
		"id2",
		stan.NatsURL(clusterUrls),
		stan.NatsOptions(
			nats.ReconnectWait(time.Second*4),
			nats.Timeout(time.Second*4)),)

	if err != nil {
		return err
	}

	_, err = sc.Subscribe(cfg.Subject, func(msg *stan.Msg) {
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

	log.Printf("Listening %s", cfg.Subject)

	return nil
}