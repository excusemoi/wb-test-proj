package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)


func main() {
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

	/*sqlStatement := `
	INSERT INTO TestTable (test)
	VALUES ('123')`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}*/

	fmt.Println("Nichuya ne slomalos")
}

