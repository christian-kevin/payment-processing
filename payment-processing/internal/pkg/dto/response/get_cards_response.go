package response

import dto "spenmo/payment-processing/payment-processing/internal/pkg/store"

type GetCards struct {
	Cards []*Card `json:"cards"`
}

type Card struct {
	CardID     int64  `json:"card_id"`
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
}

func FromCardStore(card *dto.Card) *Card {
	c := Card{
		CardID:     card.ID,
		CardNumber: card.CardNumber,
		ExpiryDate: card.ExpiryDate,
		Name:       card.Name,
		CreatedAt:  card.CreatedAt,
	}
	return &c
}
