package core_http_types

import (
	"encoding/json"

	"github.com/Kosvu/todoapp-golang/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

// Создаем метод, чтобы при декодировании, оно декодировалось по нашим правилам
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value
	return nil
}

// Метод который возвращает доменный nulable
func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}
