package main

import (
	"context"
	"flag"
	"fmt"
	"influxcluster/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
)

func main() {
	port := flag.String("p", "8090", "Public Server Port")
	fmt.Println("port:", *port)
	flag.Parse()

	r := router.NewRouter()
	e := r.PublicServer()
	defer r.Dispose()

	publicServer := &http.Server{
		Addr:    ":" + *port,
		Handler: e,
	}
	go func() {
		// service connections
		if err := publicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	timeoutSeconds := viper.GetDuration("gracefulShutdown.timeoutSeconds")

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()
	if err := publicServer.Shutdown(ctx); err != nil {
		log.Fatal("Public Server Shutdown:", err)
	}
	log.Println("Public Server exiting")
}
