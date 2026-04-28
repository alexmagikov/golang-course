package dto

type SubscriptionResponse struct {
	ID    int64  `json:"id"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}
