package endpoint

import (
    "context"
    "github.com/glitchdawg/campaign-targeting-engine/internal/models"
    "github.com/glitchdawg/campaign-targeting-engine/internal/service"

    "github.com/go-kit/kit/endpoint"
)

type GetCampaignsRequest struct {
    App     string
    Country string
    OS      string
}

type GetCampaignsResponse struct {
    Campaigns []models.DeliveryResponse `json:"campaigns,omitempty"`
    Error     string                    `json:"error,omitempty"`
}

func MakeGetCampaignsEndpoint(svc service.DeliveryService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(GetCampaignsRequest)
        campaigns, err := svc.GetCampaigns(req.App, req.Country, req.OS)
        if err != nil {
            return GetCampaignsResponse{Error: err.Error()}, nil
        }
        var resp []models.DeliveryResponse
        for _, c := range campaigns {
            resp = append(resp, models.DeliveryResponse{
                CID:   c.ID,
                Image: c.Image,
                CTA:   c.CTA,
            })
        }
        return GetCampaignsResponse{Campaigns: resp}, nil
    }
}