package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/convertutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func (m *Module) GetWallet(ctx context.Context, userID int64) (*response.GetWallet, error) {
	w, err := m.walletStore.GetWalletByUserID(ctx, m.walletStore.DBX(), userID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	if w == nil {
		return nil, errutil.ErrWalletNotFound
	}

	wallet := response.GetWallet{
		WalletID: w.ID,
		Balance:  convertutil.DeflatePrice(&w.Balance),
		Country:  w.Country,
	}

	return &wallet, nil
}
