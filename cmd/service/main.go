package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/glitchdawg/campaign-targeting-engine/internal/endpoint"
	"github.com/glitchdawg/campaign-targeting-engine/internal/service"
	"github.com/glitchdawg/campaign-targeting-engine/internal/storage"
	"github.com/glitchdawg/campaign-targeting-engine/internal/transport"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	if connStr == "" {
		log.Fatal("DB connection error")
	}

	store, err := storage.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	campaigns, err := store.GetCampaigns()
	if err != nil {
		log.Fatalf("failed to load campaigns: %v", err)
	}
	rules, err := store.GetTargetingRules()
	if err != nil {
		log.Fatalf("failed to load targeting rules: %v", err)
	}

	svc := service.NewDeliveryService(campaigns, rules)
	ep := endpoint.MakeGetCampaignsEndpoint(svc)
	handler := transport.NewHTTPHandler(ep)

	http.Handle("/v1/delivery", handler)
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
