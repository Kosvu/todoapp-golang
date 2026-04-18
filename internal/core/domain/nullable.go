package domain

/*
Структура чтобы отделять null от пустого поля в JSON

Используем дженерик чтобы не привязываться к определенному типу
*/
type Nullable[T any] struct {
	Value *T
	Set   bool
}
