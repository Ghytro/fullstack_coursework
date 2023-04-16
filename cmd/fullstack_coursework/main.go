package main

import (
	"log"
	"os"
)

func main() {
	pgUrl := os.Getenv("DB_URL")
	mongoUrl := os.Getenv("MONGO_URL")
	// todo: viper
	if err := migrate(pgUrl, "up"); err != nil {
		log.Fatal(err)
	}
	serve(pgUrl, mongoUrl)
}
