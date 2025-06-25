package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
    "time"

    "github.com/glitchdawg/campaign-targeting-engine/internal/endpoint"
    "github.com/glitchdawg/campaign-targeting-engine/internal/models"
    "github.com/glitchdawg/campaign-targeting-engine/internal/service"
    "github.com/glitchdawg/campaign-targeting-engine/internal/storage"
    "github.com/glitchdawg/campaign-targeting-engine/internal/transport"
    "github.com/joho/godotenv"
)

type HotReloadDeliveryService struct {
    svc    *service.DeliveryServiceImpl
    mutex  sync.RWMutex
}

func NewHotReloadDeliveryService(campaigns map[string]models.Campaign, rules map[string]models.TargetingRule) *HotReloadDeliveryService {
    return &HotReloadDeliveryService{
        svc: service.NewDeliveryServiceImpl(campaigns, rules),
    }
}

func (h *HotReloadDeliveryService) GetCampaigns(app, country, os string) ([]models.Campaign, error) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    return h.svc.GetCampaigns(app, country, os)
}

func (h *HotReloadDeliveryService) Reload(campaigns map[string]models.Campaign, rules map[string]models.TargetingRule) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    h.svc = service.NewDeliveryServiceImpl(campaigns, rules)
}

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
        log.Fatal("Database connection string is empty")
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

    hotSvc := NewHotReloadDeliveryService(campaigns, rules)
    ep := endpoint.MakeGetCampaignsEndpoint(hotSvc)
    handler := transport.NewHTTPHandler(ep)

    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        for range ticker.C {
            campaigns, err := store.GetCampaigns()
            if err != nil {
                log.Printf("failed to reload campaigns: %v", err)
                continue
            }
            rules, err := store.GetTargetingRules()
            if err != nil {
                log.Printf("failed to reload targeting rules: %v", err)
                continue
            }
            hotSvc.Reload(campaigns, rules)
            log.Println("Reloaded campaigns and rules from DB")
        }
    }()

    http.Handle("/v1/delivery", handler)
    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}