package models
type Campaign struct {
	ID     string `json:"cid"`
	Name   string `json:"name"`
	Image  string `json:"img"`
	CTA    string `json:"cta"`
	Status string `json:"status"`
}

