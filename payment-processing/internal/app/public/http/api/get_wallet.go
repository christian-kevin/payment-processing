package api

import (
	"context"
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/contextutil"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/errutil"
)

func GetWallet(
	handlerFunc func(ctx context.Context, userID int64) (*response.GetWallet, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := contextutil.ExtractUserID(r.Context())
		if err != nil {
			response.WriteResponse(w, nil, errutil.ErrUnauthorized)
			return
		}

		resp, err := handlerFunc(r.Context(), userID)
		if err != nil {
			response.WriteResponse(w, nil, err)
			return
		}

		response.WriteResponse(w, resp, nil)
	}
}
