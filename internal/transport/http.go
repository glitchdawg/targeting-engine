package transport

import (
    "context"
    "encoding/json"
    "net/http"
    "targeting-engine/internal/endpoint"

    kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(ep endpoint.Endpoint) http.Handler {
    return kithttp.NewServer(
        ep,
        decodeRequest,
        encodeResponse,
    )
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
    app := r.URL.Query().Get("app")
    country := r.URL.Query().Get("country")
    os := r.URL.Query().Get("os")
    return endpoint.GetCampaignsRequest{
        App:     app,
        Country: country,
        OS:      os,
    }, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    resp := response.(endpoint.GetCampaignsResponse)
    if resp.Error != "" {
        if resp.Error == "missing app param" || resp.Error == "missing country param" || resp.Error == "missing os param" {
            w.WriteHeader(http.StatusBadRequest)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        return json.NewEncoder(w).Encode(map[string]string{"error": resp.Error})
    }
    if len(resp.Campaigns) == 0 {
        w.WriteHeader(http.StatusNoContent)
        return nil
    }
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(resp.Campaigns)
}