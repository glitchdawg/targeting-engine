package service

import (
	"errors"
	"strings"

	"github.com/glitchdawg/campaign-targeting-engine/internal/models"
)

type DeliveryService interface {
    GetCampaigns(app, country, os, state string) ([]models.Campaign, error)
    SetCountryStates(map[string][]string)
}


type DeliveryServiceImpl struct {
    Campaigns      map[string]models.Campaign
    TargetingRules map[string]models.TargetingRule
    CountryStates  map[string][]string
}

func NewDeliveryServiceImpl(campaigns map[string]models.Campaign, rules map[string]models.TargetingRule) *DeliveryServiceImpl {
    return &DeliveryServiceImpl{Campaigns: campaigns, TargetingRules: rules}
}

func (s *DeliveryServiceImpl) FindParentCountry(state string) string{
    if state!="" && s.CountryStates!=nil{
        parentCountry := ""
        for c, states := range s.CountryStates {
            for _, s := range states {
                if strings.EqualFold(s, state) {
                    parentCountry = c
                    break
                }
            }
        }
        if parentCountry != "" {
            return parentCountry
        }
    }
    return ""
}
func (s *DeliveryServiceImpl) SetCountryStates(cs map[string][]string){
    s.CountryStates = cs
}
func (s *DeliveryServiceImpl) GetCampaigns(app, country, os, state string) ([]models.Campaign, error) {

    parentCountry := s.FindParentCountry(state)

    for _, rule := range s.TargetingRules {
        if len(rule.ExcludeCountry)>0 && parentCountry!=""{
            for _, exc := range rule.ExcludeCountry {
                if strings.EqualFold(exc, parentCountry) {
                    return nil, errors.New("parent Country cannot be excluded")
                }
            }
        }
    }
    if app == "" {
        return nil, errors.New("missing app param")
    }
    if country == "" {
        return nil, errors.New("missing country param")
    }
    if os == "" {
        return nil, errors.New("missing os param")
    }
    if state==""{
        return nil, errors.New("missing state param")
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
        if !matchRule(rule, app, country, os,state) {
            continue
        }
        result = append(result, camp)
    }
    return result, nil
}

func matchRule(rule models.TargetingRule, app, country, os, state string) bool {
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
    //State
    //state:country
    if len(rule.IncludeState)>0 && !in(rule.IncludeState,state){
        return false
    }
    if len(rule.ExcludeState)>0 && in(rule.ExcludeState,state){
        return false
    }

    return true
}