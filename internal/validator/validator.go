package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

/*
本package提供一组，预定义的字段校验方法，比如判断字段长度是否合适、是否为空、是否在枚举值里
可以供任意调用方使用，提供一个valid() bool 方法，判断一组字段是否满足要求，不满足时则提供一组提示信息
*/

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	FieldErrors map[string]string
}

// 如果false，则记录一组提示信息
func (v *Validator) AddFildError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFildError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// 一组内置检验
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermmittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func MinChars(chars string, n int) bool {
	return utf8.RuneCountInString(chars) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
