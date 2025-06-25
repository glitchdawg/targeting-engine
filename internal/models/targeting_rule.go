package models

import "sync"

type TargetingRule struct {
	CampaignID     string   `json:"campaign_id"`
	IncludeCountry []string `json:"include_country,omitempty"`
	ExcludeCountry []string `json:"exclude_country,omitempty"`
	IncludeOS      []string `json:"include_os,omitempty"`
	ExcludeOS      []string `json:"exclude_os,omitempty"`
	IncludeApp     []string `json:"include_app,omitempty"`
	ExcludeApp     []string `json:"exclude_app,omitempty"`
}

type TargetingEngine struct {
	campaigns map[string]*Campaign
	rules     map[string]*TargetingRule
	mutex     sync.RWMutex
}
