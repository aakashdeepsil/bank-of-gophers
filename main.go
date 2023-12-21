package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

var Seed *bool

func init() {
	Seed = flag.Bool("seed", false, "seed the database with some accounts")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Loaded environment variables successfully!!!")

}

func main() {
	log.Println("Welcome to the Bank of Gophers!!!")

	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database successfully!!!")

	if *Seed {
		log.Printf("Seeding the database with some accounts...")

		if err := seedAccounts(store); err != nil {
			log.Fatal(err)
		}
	}

	server := NewAPIServer(":8080", store)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("Goodbye!!!")
}
