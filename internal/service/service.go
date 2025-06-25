package service

import (
    "errors"
    "strings"
    "targeting-engine/internal/models"
)

type DeliveryService interface {
    GetCampaigns(app, country, os string) ([]models.Campaign, error)
}

type deliveryService struct {
    Campaigns      map[string]models.Campaign
    TargetingRules map[string]models.TargetingRule
}

func NewDeliveryService(campaigns map[string]models.Campaign, rules map[string]models.TargetingRule) DeliveryService {
    return &deliveryService{Campaigns: campaigns, TargetingRules: rules}
}

func (s *deliveryService) GetCampaigns(app, country, os string) ([]models.Campaign, error) {
    if app == "" {
        return nil, errors.New("missing app param")
    }
    if country == "" {
        return nil, errors.New("missing country param")
    }
    if os == "" {
        return nil, errors.New("missing os param")
    }

    var result []models.Campaign
    for cid, camp := range s.Campaigns {
        if strings.ToUpper(camp.Status) != "ACTIVE" {
            continue
        }
        rule, ok := s.TargetingRules[cid]
        if !ok {
            continue
        }
        if !matchRule(rule, app, country, os) {
            continue
        }
        result = append(result, camp)
    }
    return result, nil
}

func matchRule(rule models.TargetingRule, app, country, os string) bool {
    in := func(list []string, val string) bool {
        for _, v := range list {
            if strings.EqualFold(v, val) {
                return true
            }
        }
        return false
    }
    // App
    if len(rule.IncludeApp) > 0 && !in(rule.IncludeApp, app) {
        return false
    }
    if len(rule.ExcludeApp) > 0 && in(rule.ExcludeApp, app) {
        return false
    }
    // Country
    if len(rule.IncludeCountry) > 0 && !in(rule.IncludeCountry, country) {
        return false
    }
    if len(rule.ExcludeCountry) > 0 && in(rule.ExcludeCountry, country) {
        return false
    }
    // OS
    if len(rule.IncludeOS) > 0 && !in(rule.IncludeOS, os) {
        return false
    }
    if len(rule.ExcludeOS) > 0 && in(rule.ExcludeOS, os) {
        return false
    }
    return true
}