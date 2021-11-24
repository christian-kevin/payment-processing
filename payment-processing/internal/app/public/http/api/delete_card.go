package api

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/convertutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/request"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func DeleteCard(
	handlerFunc func(ctx context.Context, userID int64, cardID int64) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.DeleteCard
		if err := convertutil.MapRequestToModel(r, &req); err != nil {
			log.Get().Errorf(r.Context(), "cannot map delete card request, err: %w", err)
			response.WriteResponse(w, nil, errutil.ErrInvalidParam)
			return
		}

		userID, err := contextutil.ExtractUserID(r.Context())
		if err != nil {
			response.WriteResponse(w, nil, errutil.ErrUnauthorized)
			return
		}

		err = handlerFunc(r.Context(), userID, req.CardID)
		if err != nil {
			response.WriteResponse(w, nil, err)
			return
		}

		response.WriteResponse(w, nil, nil)
	}
}
