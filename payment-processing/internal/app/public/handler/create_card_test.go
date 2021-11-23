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

func TestCreateCard(t *testing.T) {
	log.InitializeLogger("TestCreateCard")
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
		name             string
		mockExpectations func()
		expectations     func(t *testing.T, res *response.CreateCard, err error)
	}{
		{
			caseName:         "when user don't have wallet",
			ctx:              contextWithUserID,
			userID:           int64(1),
			name:             "X",
			mockExpectations: func() {
				mockWalletStore.EXPECT().GetWalletByUserID(contextWithUserID, db, gomock.Any()).Return(nil, nil)
				mockWalletStore.EXPECT().DBX().Return(db)
			},
			expectations: func(t *testing.T, res *response.CreateCard, err error) {
				assert.Error(t, err, "should be error")
				if !errors.Is(err, errutil.ErrWalletNotFound) {
					t.Errorf("error %v is not an instance of error %v", err, errutil.ErrWalletNotFound)
				}
			},
		},
		{
			caseName:         "when can't get lock create card",
			ctx:              contextWithUserID,
			userID:           int64(1),
			name:             "X",
			mockExpectations: func() {
				mockWalletStore.EXPECT().GetWalletByUserID(contextWithUserID, db, gomock.Any()).Return(&store.Wallet{}, nil)
				mockWalletStore.EXPECT().DBX().Return(db)
				mockRedisCardStore.EXPECT().LockCreateCard(gomock.Any()).Return(errors.New("any error"))
			},
			expectations: func(t *testing.T, res *response.CreateCard, err error) {
				assert.Error(t, err, "should be error")
				if !errors.Is(err, errutil.ErrDuplicateRequest) {
					t.Errorf("error %v is not an instance of error %v", err, errutil.ErrDuplicateRequest)
				}
			},
		},
		{
			caseName:         "when success create card",
			ctx:              contextWithUserID,
			userID:           int64(1),
			name:             "X",
			mockExpectations: func() {
				mockWalletStore.EXPECT().GetWalletByUserID(contextWithUserID, db, gomock.Any()).Return(&store.Wallet{}, nil)
				mockWalletStore.EXPECT().DBX().Return(db)
				mockRedisCardStore.EXPECT().LockCreateCard(gomock.Any()).Return(nil)
				defer mockRedisCardStore.EXPECT().ReleaseLockCreateCard(gomock.Any()).Return(nil)
				mockCardStore.EXPECT().CreateCard(contextWithUserID, db, gomock.Any()).Return(int64(1), nil)
				mockCardStore.EXPECT().DBX().Return(db)
			},
			expectations: func(t *testing.T, res *response.CreateCard, err error) {
				assert.Nil(t, err)
				assert.Equal(t, res.CardID, int64(1))
			},
		},
	}

	for _, l := range tt {
		t.Log(l.caseName)

		if l.mockExpectations != nil {
			l.mockExpectations()
		}

		res, err := m.CreateCard(l.ctx, l.userID, l.name)

		l.expectations(t, res, err)
	}
}
