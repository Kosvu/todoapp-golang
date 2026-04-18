package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Kosvu/todoapp-golang/internal/core/errors"
	core_logger "github.com/Kosvu/todoapp-golang/internal/core/logger"
	"go.uber.org/zap"
)

/*
Структура HTTP ответа
содержит loger и writer
*/
type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

/*
конструктор
*/
func NewHTTPRespsonseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

/*
Метод который устанавливает статус код и энкодит тело ответа
*/
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	// пишем статус код
	h.rw.WriteHeader(statusCode)
	// энкодим тело(отправляем его)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		//логируем ошибку в случае чего
		h.log.Error("Write http response", zap.Error(err))
	}
}

// Ответ без тела (при удалении к примеру)
func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

/*
метод который мапит ошибку и в зависимости от этого
выставляет статус код и уровень логирования
*/
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		// переменная которая содержит в себе функцию такой же структуры
		// как и h.log.Warn например и т.д
		logFunc func(string, ...zap.Field)
	)

	//мапим ошибку и придаем ей нужный статус код и logFunc
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	//после чего вызываем метод errorResponse
	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// Метод который выставляет определенный статус код и панику
// после чего просто возвразает метод errorResponse
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

/*
Метод которвый выписывает ошибку в красивом формате
*/
func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	// красивый формат в виде мапы
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	// вызывает метод который пишет json
	h.JSONResponse(
		response,
		statusCode,
	)

}
