package domain

type PurchaseStatus string

const (
	PurchasePending    PurchaseStatus = "pending"
	PurchaseDispatched PurchaseStatus = "dispatched"
	PurchaseDelivered  PurchaseStatus = "delivered"
)

type Purchase struct {
	ID     string         `json:"id"`
	Status PurchaseStatus `json:"status"`
}

func (p *Purchase) ValidatePurchaseAssignment(purchase *Purchase, route *Route) error {
	validators := []func(*Purchase, *Route) error{
		validateNotDuplicate,
		validateIsValidStatus,
	}

	for _, validator := range validators {
		if err := validator(purchase, route); err != nil {
			return err
		}
	}

	return nil
}

func validateNotDuplicate(purchase *Purchase, route *Route) error {
	for _, p := range route.Purchases {
		if p.ID == purchase.ID {
			return NewError(ErrorPurchaseAlreadyExists, "purchase already exist", nil)
		}
	}
	return nil
}

func validateIsValidStatus(purchase *Purchase, route *Route) error {
	if purchase.Status != PurchasePending {
		return NewError(ErrorPurchaseNotInPendingStatus, "not in pending status", nil)
	}

	return nil
}
