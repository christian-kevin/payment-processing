package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	dto "spenmo/payment-processing/payment-processing/internal/pkg/store"
)

type limitStore struct {
	Store
	dbx *sqlx.DB
}

func NewLimitStore(db *component.Database) LimitStore {
	return &limitStore{
		Store: NewStore(db),
		dbx:   sqlx.NewDb(db.Master, "mysql"),
	}
}

const (
	insertLimitQuery = `
		insert into limits (
			parent_id,
			parent_type,
			type,
			amount
		) values (
			:parent_id,
			:parent_type,
			:type,
			:amount
	`

	getLimitsByParentIDAndParentType = `
		select
			id,
			parent_id,
			parent_type,
			type,
			amount
		from 
			limits
		where
			parent_id = ? and parent_type = ?
	`
)

func (l *limitStore) CreateLimit(ctx context.Context, execer Execer, limit *dto.Limit) (limitID int64, err error) {
	res, err := execer.NamedExecContext(ctx, insertLimitQuery, limit)
	if err != nil {
		return -1, fmt.Errorf("failed to insert limit %v: %w", limit, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (l *limitStore) GetLimits(ctx context.Context, querier Querier, parentType int, parentID int64) (res []*dto.Limit, err error) {
	var limits []dto.Limit
	res = make([]*dto.Limit, 0, 0)
	if err := querier.SelectContext(ctx, &limits, getLimitsByParentIDAndParentType, parentID, parentType); err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	for _, v := range limits {
		limit := v
		res = append(res, &limit)
	}

	return res, nil
}
