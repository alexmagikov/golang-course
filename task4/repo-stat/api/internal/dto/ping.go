package dto

import "repo-stat/api/internal/domain"

type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func FromDomainPingResult(result domain.PingResult) PingResponse {
	services := make([]ServiceStatus, len(result.Services))
	for i, s := range result.Services {
		services[i] = ServiceStatus{
			Name:   s.Name,
			Status: string(s.Status),
		}
	}

	return PingResponse{
		Status:   result.Status,
		Services: services,
	}
}
