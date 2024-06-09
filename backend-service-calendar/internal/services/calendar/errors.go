package calendar

import "errors"

var (
	//ErrEmailNotValid        = errors.New("email is not valid")
	ErrNilEntity = errors.New("entity is nil")
	//ErrUserAlreadyExist     = errors.New("user already exist")
	ErrCalendarAlreadyExist = errors.New("calendar already exist")
	ErrCalendarNotFound     = errors.New("calendar not found")
	//ErrWrongPassword        = errors.New("wrong password")
)
