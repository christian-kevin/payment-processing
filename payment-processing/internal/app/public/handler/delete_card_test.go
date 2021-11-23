package handler

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/testutils"
	"spenmo/payment-processing/payment-processing/mocks/redismocks"
	"spenmo/payment-processing/payment-processing/mocks/storemocks"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"testing"
)

func TestDeleteWallet(t *testing.T) {
	log.InitializeLogger("TestDeleteWallet")
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
		cardID           int64
		mockExpectations func()
		expectations     func(t *testing.T, err error)
	}{
		{
			caseName:         "when card not found",
			ctx:              contextWithUserID,
			userID:           int64(1),
			cardID:           int64(1),
			mockExpectations: func() {
				mockCardStore.EXPECT().GetCardByID(contextWithUserID, db, gomock.Any()).Return(nil, nil)
				mockCardStore.EXPECT().DBX().Return(db)
			},
			expectations: func(t *testing.T, err error) {
				assert.Error(t, err, "should be error")
				if !errors.Is(err, errutil.ErrCardNotFound) {
					t.Errorf("error %v is not an instance of error %v", err, errutil.ErrCardNotFound)
				}
			},
		},
	}

	for _, l := range tt {
		t.Log(l.caseName)

		if l.mockExpectations != nil {
			l.mockExpectations()
		}

		err := m.DeleteCard(l.ctx, l.userID, l.cardID)

		l.expectations(t, err)
	}
}

