package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type PingUseCase struct {
	processorClient  ProcessorClient
	subscriberClient SubscriberClient
}

func NewPing(processorClient ProcessorClient, subscriberClient SubscriberClient) *PingUseCase {
	return &PingUseCase{
		processorClient:  processorClient,
		subscriberClient: subscriberClient,
	}
}

func (u *PingUseCase) Execute(ctx context.Context) domain.PingResult {
	services := make([]domain.ServiceStatus, 0, 2)
	allUp := true

	processorStatus := u.processorClient.Ping(ctx)
	services = append(services, domain.ServiceStatus{
		Name:   "processor",
		Status: processorStatus,
	})
	if processorStatus != domain.PingStatusUp {
		allUp = false
	}

	subscriberStatus := u.subscriberClient.Ping(ctx)
	services = append(services, domain.ServiceStatus{
		Name:   "subscriber",
		Status: subscriberStatus,
	})
	if subscriberStatus != domain.PingStatusUp {
		allUp = false
	}

	status := "ok"
	if !allUp {
		status = "degraded"
	}

	return domain.PingResult{
		Status:   status,
		Services: services,
	}
}
