package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Kosvu/todoapp-golang/internal/core/errors"
)

// структура пользовател
type User struct {
	ID      int
	Version int

	FullName string
	// указатель на string, потому что поле не обязательное и может быть nil
	PhoneNumber *string
}

// конструктор
func NewUser(
	id int,
	version int,
	full_name string,
	phone_number *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
}

// конструктор если не заданы поля id и version
func NewUserUninitialized(
	full_name string,
	phone_number *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		full_name,
		phone_number,
	)

}

// Валидация пользователя
func (u *User) Validate() error {
	fullNameLenght := len([]rune(u.FullName))

	if fullNameLenght < 3 || fullNameLenght > 100 {
		return fmt.Errorf(
			"invalid `FullName` len %d: %w",
			fullNameLenght,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))

		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"ivalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

// Валидация патча
func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("`FullName` cant't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
