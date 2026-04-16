package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/Kosvu/todoapp-golang/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

// Валидатор запроса
var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {

	//сначала все декодируем в нашу переменную из тела запроса
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var (
		err error
	)

	//проверяем имеет ли dest метод validate, чтобы понять если правила на поля
	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		//если закодировали, теперь проверяем на валидность
		// (валидности прописаны в структуре самой фичи)
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
