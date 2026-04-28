package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type ServiceStatus struct {
	Name   string
	Status PingStatus
}

type PingResult struct {
	Status   string
	Services []ServiceStatus
}
