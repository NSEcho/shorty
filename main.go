package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lateralusd/shorty/db"
	"github.com/lateralusd/shorty/handler"
)

func main() {
	// Command line flags
	port := flag.Int("port", 8080, "port number to bind web server to")
	dbName := flag.String("db", "links.db", "the name of the database")
	timeout := flag.Int("timeout", 1, "timeout for database")
	flag.Parse()

	// Create functional options
	bucketName := withBucketName(*dbName)
	timeoutVal := withTimeout(*timeout)

	addr := fmt.Sprintf(":%d", *port)

	// Initialize database with our functional options
	var err error
	database, err := db.InitDatabase(bucketName, timeoutVal)
	if err != nil {
		log.Fatal("Coult not initialize database", err)
	}

	env := &handler.Env{
		DB: database,
	}

	http.Handle("/", handler.Handler{env, handler.IndexPath})
	http.Handle("/shorty", handler.Handler{env, handler.ShortyPath})

	fmt.Printf("[*] Starting the server on port %d\n", *port)
	http.ListenAndServe(addr, nil)
}

func withBucketName(bucketName string) db.ConfigOption {
	return func(cfg *db.Config) {
		cfg.Bucket = []byte(bucketName)
	}
}

func withTimeout(timeout int) db.ConfigOption {
	return func(cfg *db.Config) {
		cfg.Timeout = timeout
	}
}
