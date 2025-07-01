package models

import (
	"errors"
)

// 自定义error信息
var (
	ErrNoRecord           = errors.New("models: no matching record found!")
	ErrInvalidCredentials = errors.New("models: invalid credentials!")
	ErrDuplicateEmail     = errors.New("models: duplicate email!")
)
