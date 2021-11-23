package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/cardutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/store"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func (m *Module) CreateCard(ctx context.Context, userID int64, name string) (*response.CreateCard, error) {
	w, err := m.walletStore.GetWalletByUserID(ctx, m.walletStore.DBX(), userID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	if w == nil {
		return nil, errutil.ErrWalletNotFound
	}

	err = m.rCardStore.LockCreateCard(w.ID)
	if err != nil {
		log.Get().Errorf(ctx, "failed to get lock from redis, err: %s", err.Error())
		return nil, errutil.ErrDuplicateRequest
	}

	defer m.rCardStore.ReleaseLockCreateCard(w.ID)

	cardNumber := cardutil.GenerateCardNumber()
	expiryDate := cardutil.GenerateExpiryDate(w.Country)

	card := store.Card{
		WalletID:   w.ID,
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		Name:       name,
	}

	id, err := m.cardStore.CreateCard(ctx, m.cardStore.DBX(), &card)
	if err != nil {
		log.Get().Errorf(ctx, "failed to get create card, err: %s", err.Error())
		return nil, errutil.ErrDuplicateRequest
	}

	res := response.CreateCard{
		CardID:     id,
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		Name:       name,
	}

	return &res, nil
}


