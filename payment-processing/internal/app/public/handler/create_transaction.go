package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	"spenmo/payment-processing/payment-processing/internal/pkg/convertutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/retry_util"
	"spenmo/payment-processing/payment-processing/internal/pkg/store"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

type limitToRetry struct {
	limit int64
	retType int
}

func (m *Module) CreateTransaction(ctx context.Context, value string, cardNumber, expiryDate string) error {
	card, err := m.cardStore.GetCardByNumberAndExpiryDate(ctx, m.cardStore.DBX(), cardNumber, expiryDate)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch card, err: %s", err.Error())
		return errutil.ErrServerError
	}

	if card == nil {
		return errutil.ErrCardNotFound
	}

	limits, err := m.limitStore.GetLimits(ctx, m.limitStore.DBX(), constant.ParentLimitTypeCard, card.ID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch card limit, err: %s", err.Error())
		return errutil.ErrServerError
	}

	valInflated, err := convertutil.InflatePrice(value)
	if err != nil {
		log.Get().Errorf(ctx, "error inflate price, err: %s", err.Error())
		return errutil.ErrInvalidParam
	}

	for _, l := range limits {
		if *valInflated > l.Amount {
			return errutil.ErrCardLimitExceeded
		}
	}

	return m.processTransaction(ctx, *valInflated, limits, card)
}

func (m *Module) processTransaction(ctx context.Context, valInflated int64, limits []*store.Limit, card *store.Card) (err error) {
	err = m.rWalletStore.LockUpdateWallet(card.WalletID)
	if err != nil {
		log.Get().Errorf(ctx, "error get lock, err: %s", err.Error())
		return errutil.ErrDuplicateRequest
	}
	defer m.rWalletStore.ReleaseLockUpdateWallet(card.WalletID)

	tx, err := m.walletStore.BeginX()
	if err != nil {
		log.Get().Errorf(ctx, "error get tx, err: %s", err.Error())
		return errutil.ErrServerError
	}

	w, err := m.walletStore.GetWalletByID(ctx, tx, card.WalletID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return errutil.ErrServerError
	}

	w.Balance-=valInflated
	if w.Balance <= 0 {
		return errutil.ErrNotEnoughBalance
	}

	isSuccess := make([]*limitToRetry, 0, 0)
	for _, l := range limits {
		err = m.rCardStore.CreateTransaction(l.Amount, valInflated, card.ID, l.Type)
		if err == nil {
			lt := limitToRetry{
				limit:   l.Amount,
				retType: l.Type,
			}
			isSuccess = append(isSuccess, &lt)
		}
	}
	defer func(err error) {
		if err != nil {
			for _, lt := range isSuccess {
				go m.rollbackLimit(ctx, valInflated, lt.limit, card.ID, lt.retType)
			}
		}
	}(err)

	if len(isSuccess) != len(limits) {
		err = errutil.ErrCardLimitExceeded
		return err
	}

	err = m.walletStore.ModifyBalance(ctx, tx, w.Balance, w.ID)
	if err != nil {
		log.Get().Errorf(ctx, "error modify wallet balance, err: %s", err.Error())
		return errutil.ErrServerError
	}

	wL := store.WalletBalanceLog{
		WalletID:  w.ID,
		Amount:    -1*valInflated,
	}
	err = m.walletBalanceLogStore.CreateLog(ctx, tx, &wL)
	if err != nil {
		log.Get().Errorf(ctx, "error create wallet log, err: %s", err.Error())
		return errutil.ErrServerError
	}

	cL := store.CardTransactionLog{
		CardID:    card.ID,
		Amount:    valInflated,
	}

	err = m.cardTransactionLogStore.CreateLog(ctx, tx, &cL)
	if err != nil {
		log.Get().Errorf(ctx, "error create card log, err: %s", err.Error())
		return errutil.ErrServerError
	}

	err = m.walletStore.CommitX(tx)
	if err != nil {
		log.Get().Errorf(ctx, "error failed to commit tx, err: %s", err.Error())
		return errutil.ErrServerError
	}

	return nil
}

func (m *Module) rollbackLimit(ctx context.Context, valInflated int64, limit int64, cardID int64, retType int) {
	funcToRetry := func() error {
		return m.rCardStore.CreateTransaction(limit, -1*valInflated, cardID, retType)
	}

	errRetry := retry_util.ExponentialRetry(ctx, 5, funcToRetry)
	if errRetry != nil {
		log.Get().Errorf(ctx, "failed to rollback limit. err %s, cardID %d, retType: %d", errRetry.Error(), cardID, retType)
	}
}
