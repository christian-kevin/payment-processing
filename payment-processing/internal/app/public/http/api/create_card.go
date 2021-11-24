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

func CreateCard(
	handlerFunc func(ctx context.Context, userID int64, name string) (*response.CreateCard, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.CreateCard
		if err := convertutil.MapRequestToModel(r, &req); err != nil {
			log.Get().Errorf(r.Context(), "cannot map create card request, err: %w", err)
			response.WriteResponse(w, nil, errutil.ErrInvalidParam)
			return
		}

		userID, err := contextutil.ExtractUserID(r.Context())
		if err != nil {
			response.WriteResponse(w, nil, errutil.ErrUnauthorized)
			return
		}

		resp, err := handlerFunc(r.Context(), userID, req.Name)
		if err != nil {
			response.WriteResponse(w, nil, err)
			return
		}

		response.WriteResponse(w, resp, nil)
	}
}
