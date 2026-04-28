package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type Ping struct{}

func NewPing() *Ping {
	return &Ping{}
}

func (u *Ping) Ping(ctx context.Context) (domain.PingStatus, error) {
	return domain.PingStatusUp, nil
}
