package repositories

import (
	"context"

	"github.com/uMakeMeCrazy/fravega_tech/pkg/logger"
	"go.uber.org/zap"

	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"
)

type PurchasesRepo struct {
}

func NewPurchasesRepo() *PurchasesRepo {
	return &PurchasesRepo{}
}

func (p *PurchasesRepo) FindByID(ctx context.Context, purchaseID string) (*domain.Purchase, error) {
	return &domain.Purchase{
		ID:     purchaseID,
		Status: domain.PurchasePending,
	}, nil
}

func (p *PurchasesRepo) SendEmailNotification(ctx context.Context, purchaseID string) error {
	logger.Info(ctx, "send email notification for purchase delivered", zap.String(domain.PurchaseID, purchaseID))

	return nil
}
