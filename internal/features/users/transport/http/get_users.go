package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_http_request "github.com/Kosvu/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/Kosvu/todoapp-golang/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_response.NewHTTPRespsonseHandler(log, w)

	limit, offset, err := getLimifOffsetQueryParams(r)
	if err != nil {
		responceHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query params",
		)

		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responceHandler.ErrorResponse(
			err,
			"failed to get users",
		)

		return
	}

	response := GetUsersResponse(userDTOFromDomains(userDomains))

	responceHandler.JSONResponse(response, http.StatusOK)
}

func getLimifOffsetQueryParams(r *http.Request) (*int, *int, error) {

	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", limitQueryParamKey)
	}

	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", offsetQueryParamKey)
	}

	return limit, offset, err
}
