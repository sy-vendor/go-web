package main

import (
	go_web "go-web"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	go_web.Migrate()

	server, err := Create()
	if err != nil {
		panic(err)
	}

	server.Start()
	server.AwaitSignal()
}
