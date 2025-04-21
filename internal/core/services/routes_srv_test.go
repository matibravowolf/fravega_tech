package services

import (
	"context"
	"testing"
	"time"

	"github.com/uMakeMeCrazy/fravega_tech/internal/core/domain"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/uMakeMeCrazy/fravega_tech/internal/core/services/mocks"
)

func Test_CreateRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		driver  string
		vehicle string
	}

	tests := []struct {
		name          string
		args          args
		mocks         func(*mocks.MockRoutesRepo, *mocks.MockPurchasesRepo)
		expectedError error
	}{
		{
			name: "successfully CreateRoute()",
			args: args{
				driver:  "Hernan Cattaneo",
				vehicle: "Mercedez Benz",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				routesRepo.EXPECT().
					Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Route{})).
					Return(&domain.Route{
						ID:        "random_id",
						Vehicle:   "Mercedez Benz",
						Driver:    "Hernan Cattaneo",
						Status:    domain.RouteStatusPending,
						Purchases: nil,
						CreatedAt: time.Time{},
					}, nil)
			},
			expectedError: nil,
		},
		{
			name: "error when saving route fails",
			args: args{
				driver:  "Julian Perez",
				vehicle: "Camion XXL",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, _ *mocks.MockPurchasesRepo) {
				routesRepo.EXPECT().
					Save(gomock.Any(), gomock.AssignableToTypeOf(&domain.Route{})).
					Return(nil, domain.NewError(domain.ErrorUnexpected, "unexpected_error", nil))
			},
			expectedError: domain.NewError(domain.ErrorUnexpected, "unexpected_error", nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Inject mocks
			mockRouteRepo := mocks.NewMockRoutesRepo(ctrl)
			mockPurchaseRepo := mocks.NewMockPurchasesRepo(ctrl)

			if tc.mocks != nil {
				tc.mocks(mockRouteRepo, mockPurchaseRepo)
			}

			srv := NewRoutesSrv(mockRouteRepo, mockPurchaseRepo)

			// Execute
			route, err := srv.CreateRoutes(context.Background(), tc.args.driver, tc.args.vehicle)

			// Verify
			if tc.expectedError != nil {
				assert.Nil(t, route)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NotNil(t, route)
				assert.NoError(t, err)
				assert.NotEmpty(t, route.ID)
				assert.Equal(t, tc.args.driver, route.Driver)
				assert.Equal(t, tc.args.vehicle, route.Vehicle)
			}
		})
	}
}

func Test_AddPurchaseToRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		routeID    string
		purchaseID string
	}

	tests := []struct {
		name          string
		args          args
		mocks         func(*mocks.MockRoutesRepo, *mocks.MockPurchasesRepo)
		expectedError error
	}{
		{
			name: "successfully AddPurchaseToRoute()",
			args: args{
				routeID:    "route-1",
				purchaseID: "purchase-1",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				route := domain.NewRoute("vehicle-1", "driver-1")
				purchase := &domain.Purchase{ID: "purchase-1", Status: domain.PurchasePending}

				routesRepo.EXPECT().FindByID(gomock.Any(), "route-1").Return(route, nil)
				purchasesRepo.EXPECT().FindByID(gomock.Any(), "purchase-1").Return(purchase, nil)
				routesRepo.EXPECT().Save(gomock.Any(), route).Return(nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "error when route is not found",
			args: args{
				routeID:    "invalid-route",
				purchaseID: "purchase-1",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				routesRepo.EXPECT().
					FindByID(gomock.Any(), "invalid-route").
					Return(nil, domain.NewError(domain.ErrorRouteNotFound, "route not found", nil))
			},
			expectedError: domain.NewError(domain.ErrorRouteNotFound, "route not found", nil),
		},
		{
			name: "error when route is not in pending status",
			args: args{
				routeID:    "route-2",
				purchaseID: "purchase-2",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				route := domain.NewRoute("vehicle-2", "driver-2")
				route.Status = domain.RouteStatusOnWay

				routesRepo.EXPECT().FindByID(gomock.Any(), "route-2").Return(route, nil)
			},
			expectedError: domain.NewError(domain.ErrorPurchaseNotInPendingStatus, "can't add purchase to route, is not in pending status", nil),
		},
		{
			name: "error when purchase is not found",
			args: args{
				routeID:    "route-3",
				purchaseID: "invalid-purchase",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				route := domain.NewRoute("vehicle-3", "driver-3")

				routesRepo.EXPECT().FindByID(gomock.Any(), "route-3").Return(route, nil)
				purchasesRepo.EXPECT().
					FindByID(gomock.Any(), "invalid-purchase").
					Return(nil, domain.NewError(domain.ErrorPurchaseNotExist, "invalid_purchase", nil))
			},
			expectedError: domain.NewError(domain.ErrorPurchaseNotExist, "invalid_purchase", nil),
		},
		{
			name: "error when adding purchase to route fails (duplicate)",
			args: args{
				routeID:    "route-4",
				purchaseID: "purchase-4",
			},
			mocks: func(routesRepo *mocks.MockRoutesRepo, purchasesRepo *mocks.MockPurchasesRepo) {
				route := domain.NewRoute("vehicle-4", "driver-4")
				duplicate := &domain.Purchase{ID: "purchase-4", Status: domain.PurchasePending}
				_ = route.AddPurchase(duplicate) // previously added

				routesRepo.EXPECT().FindByID(gomock.Any(), "route-4").Return(route, nil)
				purchasesRepo.EXPECT().FindByID(gomock.Any(), "purchase-4").Return(duplicate, nil)
			},
			expectedError: domain.NewError(domain.ErrorPurchaseAlreadyExists, "purchase already exist", nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Inject mocks
			mockRouteRepo := mocks.NewMockRoutesRepo(ctrl)
			mockPurchaseRepo := mocks.NewMockPurchasesRepo(ctrl)

			if tc.mocks != nil {
				tc.mocks(mockRouteRepo, mockPurchaseRepo)
			}

			srv := NewRoutesSrv(mockRouteRepo, mockPurchaseRepo)

			// Execute
			err := srv.AddPurchaseToRoute(context.Background(), tc.args.routeID, tc.args.purchaseID)

			// Verify
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
