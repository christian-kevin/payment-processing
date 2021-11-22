package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	"spenmo/payment-processing/payment-processing/internal/pkg/constant"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
	"spenmo/payment-processing/payment-processing/internal/pkg/timeutil"
)

type cardStore struct {
	Store
	dbx *sqlx.DB
}

func NewCardStore(db *component.Database) CardStore {
	return &cardStore{
		Store: NewStore(db),
		dbx:   sqlx.NewDb(db.Master, "mysql"),
	}
}

const (
	insertCardQuery = `
		insert into card(
			wallet_id,
			card_number,
			expiry_date,
			name,
			created_at
		) values (
			:wallet_id,
			:card_number,
			:expiry_date,
			:name,
			:created_at
		)
	`
	queryGetCardsByWalletID = `
		select 
			id,
			wallet_id,
			card_number,
			expiry_date,
			name,
			created_at,
			is_deleted
		from card
		where
			wallet_id = ?
	`

	queryGetCardByCardNumberAndExpiryDate = `
		select 
			id,
			wallet_id,
			card_number,
			expiry_date,
			name,
			created_at,
			is_deleted
		from card
		where
			card_number = ? and expiry_date = ?
	`

	querySoftDeleteCardByID = `
		update 
			card
		set
			is_deleted = ?
		where
			id = ?
	`
)

func (c *cardStore) CreateCard(ctx context.Context, execer Execer, card *dto.Card) (cardID int64, err error) {
	card.CreatedAt = timeutil.NowMillis()

	res, err := execer.NamedExecContext(ctx, insertCardQuery, card)
	if err != nil {
		return -1, fmt.Errorf("failed to insert card %v: %w", card, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (c *cardStore) GetCards(ctx context.Context, querier Querier, walletID int64) (res []*dto.Card, err error) {
	var cards []dto.Card
	res = make([]*dto.Card, 0, 0)
	if err := querier.SelectContext(ctx, &cards, queryGetCardsByWalletID, walletID); err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	for _, v := range cards {
		card := v
		res = append(res, &card)
	}

	return res, nil
}

func (c *cardStore) GetCardByNumberAndExpiryDate(ctx context.Context, querier Querier, cardNumber string,
	expiryDate string) (*dto.Card, error) {
	var card dto.Card
	if err := querier.GetContext(ctx, &card, queryGetCardByCardNumberAndExpiryDate, cardNumber, expiryDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find card: %w", err)
	}

	return &card, nil
}

func (c *cardStore) DeleteCard(ctx context.Context, execer Execer, cardID int64) error {
	if _, err := execer.ExecContext(ctx, querySoftDeleteCardByID, constant.IsDeleted, cardID); err != nil {
		return fmt.Errorf("failed to delete card for id %d: %w", cardID, err)
	}

	return nil
}