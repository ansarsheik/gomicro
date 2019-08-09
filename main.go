package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gamegos/jsend"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func Home(w http.ResponseWriter, r *http.Request) {
	jsend.Wrap(w).
		Status(201).
		Message("ticket data found").
		Data("fake data").
		Send()
}

func main() {
	flag.Parse()

	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Println(e.Error())
	}

	r := mux.NewRouter()

	// end points
	r.HandleFunc("/", Home).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + os.Getenv("WEBSERVER_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Webserver at port " + os.Getenv("WEBSERVER_PORT"))
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}
