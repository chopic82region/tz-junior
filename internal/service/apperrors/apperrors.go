package apperrors

import "errors"

var (
	NilField   = errors.New("не заполнено обязательное поле")
	InvalidID  = errors.New("некорректный id")
	NotFound   = errors.New("запись не найдена")
	DBError    = errors.New("ошибка базы данных")
	Duplicate  = errors.New("дубликат записи")
	BadPayload = errors.New("некорректные данные")
)
