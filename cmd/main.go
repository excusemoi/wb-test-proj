package main

import (
	_ "github.com/lib/pq"
	"goProj/app"
	"log"
)

func run() error {
	a, err := app.InitApp()
	if err != nil {
		return err
	}
	if err = a.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

