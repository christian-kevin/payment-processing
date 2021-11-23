package handler

import (
	"context"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func (m *Module) GetCards(ctx context.Context, userID int64) (*response.GetCards, error) {
	w, err := m.walletStore.GetWalletByUserID(ctx, m.walletStore.DBX(), userID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch wallet, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	cards, err := m.cardStore.GetCards(ctx, m.cardStore.DBX(), w.ID)
	if err != nil {
		log.Get().Errorf(ctx, "error fetch cards, err: %s", err.Error())
		return nil, errutil.ErrServerError
	}

	res := response.GetCards{}
	r := make([]*response.Card, 0, 0)
	for _, c := range cards {
		card := response.FromCardStore(c)
		r = append(r, card)
	}

	res.Cards = r
	return &res, nil
}
