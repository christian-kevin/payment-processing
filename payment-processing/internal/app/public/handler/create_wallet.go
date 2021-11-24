package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/convertutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/store"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func (m *Module) CreateWallet(ctx context.Context, userID int64, country string) (*response.GetWallet, error) {
	err := m.rWalletStore.LockCreateWallet(userID)
	if err != nil {
		log.Get().Errorf(ctx, "failed to get lock from redis, err: %s", err.Error())
		return nil, errutil.ErrDuplicateRequest
	}

	defer m.rWalletStore.ReleaseLockCreateWallet(userID)

	tx, err := m.walletStore.BeginX()
	if err != nil {
		log.Get().Errorf(ctx, "failed to start transaction, error: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	w, err := m.walletStore.GetWalletByUserID(ctx, tx, userID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	if w != nil {
		return nil, errutil.ErrWalletAlreadyExist
	}

	balance, err := convertutil.InflatePrice("10000000")
	if err != nil {
		log.Get().Errorf(ctx, "error inflate balance, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	wallet := store.Wallet{
		Balance:    *balance,
		ParentID:   userID,
		ParentType: constant.ParentTypeUser,
		Country:    constant.CountryID,
	}

	id, err := m.walletStore.CreateWallet(ctx, tx, &wallet)
	if err != nil {
		log.Get().Errorf(ctx, "error create wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	err = m.walletStore.CommitX(tx)
	if err != nil {
		log.Get().Errorf(ctx, "error commit create wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}
	resp := response.GetWallet{
		WalletID: id,
		Balance:  convertutil.DeflatePrice(&wallet.Balance),
		Country:  country,
	}
	return &resp, nil
}
