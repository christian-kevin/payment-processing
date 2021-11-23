package handler

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/store"
	"spenmo/payment-processing/payment-processing/internal/pkg/testutils"
	"spenmo/payment-processing/payment-processing/mocks/redismocks"
	"spenmo/payment-processing/payment-processing/mocks/storemocks"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	log.InitializeLogger("TestCreateWallet")
	db, _ := testutils.PrepareMySQLMock(t)

	ctrl := gomock.NewController(t)
	mockRedisWalletStore := redismocks.NewMockRWalletStore(ctrl)
	mockRedisCardStore := redismocks.NewMockRCardStore(ctrl)
	mockWalletStore := storemocks.NewMockWalletStore(ctrl)
	mockCardStore := storemocks.NewMockCardStore(ctrl)
	mockWalletBalanceLogStore := storemocks.NewMockWalletBalanceLogStore(ctrl)
	mockLimitStore := storemocks.NewMockLimitStore(ctrl)
	mockCardTransactionLogStore := storemocks.NewMockCardTransactionLogStore(ctrl)

	m := &Module{
		rCardStore:              mockRedisCardStore,
		rWalletStore:            mockRedisWalletStore,
		walletStore:             mockWalletStore,
		cardStore:               mockCardStore,
		limitStore:              mockLimitStore,
		cardTransactionLogStore: mockCardTransactionLogStore,
		walletBalanceLogStore:   mockWalletBalanceLogStore,
	}

	contextWithUserID := context.WithValue(context.Background(), contextutil.XUserIDKey, int64(1))

	tt := []struct {
		caseName         string
		ctx              context.Context
		userID           int64
		country          string
		mockExpectations func()
		expectations     func(t *testing.T, res *response.GetWallet, err error)
	}{
		{
			caseName:         "when failed to acquire lock",
			ctx:              contextWithUserID,
			userID:           int64(1),
			country:          "id",
			mockExpectations: func() {
				mockRedisWalletStore.EXPECT().LockCreateWallet(gomock.Any()).Return(errors.New("any error"))
			},
			expectations: func(t *testing.T, res *response.GetWallet, err error) {
				assert.Error(t, err, "should be error")
				if !errors.Is(err, errutil.ErrDuplicateRequest) {
					t.Errorf("error %v is not an instance of error %v", err, errutil.ErrDuplicateRequest)
				}
			},
		},
		{
			caseName:         "when user already have wallet",
			ctx:              contextWithUserID,
			userID:           int64(1),
			country:          "id",
			mockExpectations: func() {
				mockRedisWalletStore.EXPECT().LockCreateWallet(gomock.Any()).Return(nil)
				defer mockRedisWalletStore.EXPECT().ReleaseLockCreateWallet(gomock.Any()).Return(nil)
				tx, _ := db.Beginx()

				mockWalletStore.EXPECT().BeginX().Return(tx, nil)
				mockWalletStore.EXPECT().GetWalletByUserID(contextWithUserID, tx, int64(1)).Return(&store.Wallet{}, nil)
			},
			expectations: func(t *testing.T, res *response.GetWallet, err error) {
				assert.Error(t, err, "should be error")
				if !errors.Is(err, errutil.ErrWalletAlreadyExist) {
					t.Errorf("error %v is not an instance of error %v", err, errutil.ErrWalletAlreadyExist)
				}
			},
		},
		{
			caseName:         "when success create wallet",
			ctx:              contextWithUserID,
			userID:           int64(1),
			country:          "id",
			mockExpectations: func() {
				mockRedisWalletStore.EXPECT().LockCreateWallet(gomock.Any()).Return(nil)
				defer mockRedisWalletStore.EXPECT().ReleaseLockCreateWallet(gomock.Any()).Return(nil)
				tx, _ := db.Beginx()

				mockWalletStore.EXPECT().BeginX().Return(tx, nil)
				mockWalletStore.EXPECT().GetWalletByUserID(contextWithUserID, tx, int64(1)).Return(nil, nil)
				mockWalletStore.EXPECT().CreateWallet(contextWithUserID, tx, gomock.Any()).Return(int64(1), nil)
				mockWalletStore.EXPECT().CommitX(tx).Return(nil)
			},
			expectations: func(t *testing.T, res *response.GetWallet, err error) {
				assert.Nil(t, err)
				assert.Equal(t, res.WalletID, int64(1))
			},
		},
	}

	for _, l := range tt {
		t.Log(l.caseName)

		if l.mockExpectations != nil {
			l.mockExpectations()
		}

		res, err := m.CreateWallet(l.ctx, l.userID, l.country)

		l.expectations(t, res, err)
	}
}
