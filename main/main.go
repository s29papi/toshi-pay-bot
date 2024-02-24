package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/s29papi/wag3r-bot/service"
	"github.com/s29papi/wag3r-bot/service/utils"
)

var dev = os.Args[1]

// memory_db
// api
func main() {
	if dev == "true" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	// startBot()
	startServer()
}

var id = 318902

// func startBot() {
// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
// 	bot := worker.NewWorker(signalChan)
// 	go bot.Start()
// }

func startServer() {
	a := service.AuthHandler{
		KeyFunc:     utils.KeyFunc,
		HttpHandler: service.Mux,
	}
	server := http.Server{
		Addr:    ":8080",
		Handler: a,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// https://dashboard.render.com/web/srv-cn9tbvuv3ddc73d88a20/settings
