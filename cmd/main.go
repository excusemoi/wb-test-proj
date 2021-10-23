package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"goProj/app"
	"log"
	"os"
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
	fmt.Println(os.Getwd())
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

