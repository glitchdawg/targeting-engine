package main

import (
    "log"
    "net/http"
    "targeting-engine/internal/endpoint"
    "targeting-engine/internal/models"
    "targeting-engine/internal/service"
    "targeting-engine/internal/transport"
)

func main() {
    // Example in-memory data
    campaigns := map[string]models.Campaign{
        "spotify":      {ID: "spotify", Name: "Spotify", Image: "https://somelink", CTA: "Download", Status: "ACTIVE"},
        "duolingo":     {ID: "duolingo", Name: "Duolingo", Image: "https://somelink2", CTA: "Install", Status: "ACTIVE"},
        "subwaysurfer": {ID: "subwaysurfer", Name: "Subway Surfer", Image: "https://somelink3", CTA: "Play", Status: "ACTIVE"},
    }
    rules := map[string]models.TargetingRule{
        "spotify": {
            CampaignID:     "spotify",
            IncludeCountry: []string{"US", "Canada"},
        },
        "duolingo": {
            CampaignID:  "duolingo",
            IncludeOS:   []string{"Android", "iOS"},
            ExcludeCountry: []string{"US"},
        },
        "subwaysurfer": {
            CampaignID: "subwaysurfer",
            IncludeOS:  []string{"Android"},
            IncludeApp: []string{"com.gametion.ludokinggame"},
        },
    }

    svc := service.NewDeliveryService(campaigns, rules)
    ep := endpoint.MakeGetCampaignsEndpoint(svc)
    handler := transport.NewHTTPHandler(ep)

    http.Handle("/v1/delivery", handler)
    log.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}