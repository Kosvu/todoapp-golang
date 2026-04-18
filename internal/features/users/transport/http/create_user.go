package users_transport_http

import (
	"net/http"

	"github.com/Kosvu/todoapp-golang/internal/core/domain"
	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	core_http_request "github.com/Kosvu/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/Kosvu/todoapp-golang/internal/core/transport/http/response"
)

// создание структуры на запрос
// прописаны валидации
type CreateUserRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	//omitempty значит только в случае если PhoneNumber есть
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

// создание структуры на ответ
// создаем алиас чтобы не менять название
type CreateUserResponse UserDTOResponse

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// задаем контекст запроса
	ctx := r.Context()

	// получаем логер из запроса
	log := core_logger.FromContext(ctx)

	responseHandler := core_http_response.NewHTTPRespsonseHandler(log, w)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
