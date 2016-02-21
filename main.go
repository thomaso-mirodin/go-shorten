package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/goamz/aws"
	"github.com/thomaso-mirodin/go-shorten/handlers"
	"github.com/thomaso-mirodin/go-shorten/storage"
)

func main() {
	store, err := storage.NewInmem()
	if err != nil {
		log.Fatalf("Failed to create inmem storage because '%s'", err)

	}

	n := negroni.Classic()

	r := httprouter.New()
	r.GET("/", handlers.Index)

	r.GET("/:short", handlers.GetShortHandler(store))
	r.HEAD("/:short", handlers.GetShortHandler(store))

	r.POST("/", handlers.SetShortHandler(store))
	r.PUT("/", handlers.SetShortHandler(store))
	r.POST("/:short", handlers.SetShortHandler(store))
	r.PUT("/:short", handlers.SetShortHandler(store))

	n.UseHandler(r)

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = http.ListenAndServe(net.JoinHostPort(host, port), n)
	if err != nil {
		log.Fatal(err)
	}
}
