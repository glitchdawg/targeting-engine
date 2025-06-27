package models
type DeliveryRequest struct {
	App     string `json:"app"`
	Country string `json:"country"`
	OS      string `json:"os"`
	State  string `json:"state"`
}

type DeliveryResponse struct {
	CID   string `json:"cid"`
	Image string `json:"img"`
	CTA   string `json:"cta"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}