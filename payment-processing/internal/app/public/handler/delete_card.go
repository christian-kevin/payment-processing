package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func (m *Module) DeleteCard(ctx context.Context, userID int64, cardID int64) error {
	card, err := m.cardStore.GetCardByID(ctx, m.cardStore.DBX(), cardID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch card, err: %s", err.Error())
		return errutil.ErrServerError
	}

	if card == nil {
		return errutil.ErrCardNotFound
	}

	w, err := m.walletStore.GetWalletByUserID(ctx, m.walletStore.DBX(), userID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return errutil.ErrServerError
	}

	if w == nil {
		return errutil.ErrWalletNotFound
	}

	if w.ID != card.WalletID {
		return errutil.ErrUnauthorized
	}

	err = m.cardStore.DeleteCard(ctx, m.cardStore.DBX(), cardID)
	if err != nil {
		log.Get().Errorf(ctx, "error failed to delete card, err: %s", err.Error())
		return errutil.ErrServerError
	}

	return nil
}
