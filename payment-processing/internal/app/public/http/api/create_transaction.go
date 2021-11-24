package api

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/convertutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/request"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func CreateTransaction(
	handlerFunc func(ctx context.Context, value string, cardNumber, expiryDate string) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.CreateTransaction
		if err := convertutil.MapRequestToModel(r, &req); err != nil {
			log.Get().Errorf(r.Context(), "cannot map create transaction request, err: %w", err)
			response.WriteResponse(w, nil, errutil.ErrInvalidParam)
			return
		}

		err := handlerFunc(r.Context(), req.Amount, req.CardNumber, req.ExpiryDate)
		if err != nil {
			response.WriteResponse(w, nil, err)
			return
		}

		response.WriteResponse(w, nil, nil)
	}
}
