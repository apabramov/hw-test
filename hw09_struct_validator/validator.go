package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrLen    = errors.New("error len")
	ErrRegExp = errors.New("error regexp")
	ErrMin    = errors.New("error min")
	ErrMax    = errors.New("error max")
	ErrIn     = errors.New("error in")

	ErrValidate = errors.New("error validate")
)

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for _, k := range v {
		res.WriteString(fmt.Sprintf("%s: %v ", k.Field, k.Err))
	}
	return res.String()
}

func Validate(v interface{}) error {
	var ve ValidationErrors
	if reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return ve
	}
	t := reflect.ValueOf(v).Elem()
	typ := t.Type()
	for i, f := range reflect.VisibleFields(typ) {
		if !f.IsExported() {
			continue
		}
		if tagValue, ok := typ.Field(i).Tag.Lookup("validate"); ok {
			name := typ.Field(i).Name
			switch val := t.Field(i).Interface().(type) {
			case string:
				if err := validateString(val, tagValue, name); len(err) > 0 {
					ve = append(ve, err...)
				}
			case []string:
				if err := validateArrString(val, tagValue, name); len(err) > 0 {
					ve = append(ve, err...)
				}
			case int:
				if err := validateInt(val, tagValue, name); len(err) > 0 {
					ve = append(ve, err...)
				}
			case []int:
				if err := validateArrInt(val, tagValue, name); len(err) > 0 {
					ve = append(ve, err...)
				}
			}
		}
	}
	return ve
}

func containsArray(s []string, val int) (bool, error) {
	for _, v := range s {
		if l, err := strconv.Atoi(v); err != nil {
			return false, err
		} else if l == val {
			return true, nil
		}
	}
	return false, nil
}

func checkArray(s []string, name string) error {
	if len(s) != 2 {
		return fmt.Errorf("field %s: %w", name, ErrValidate)
	}
	if s[1] == "" {
		return fmt.Errorf("field %s: %w", name, ErrValidate)
	}
	return nil
}

func validateString(val string, tagValue string, name string) ValidationErrors {
	var ve ValidationErrors
	if val == "" {
		return ve
	}
	s := strings.Split(tagValue, "|")
	for _, v := range s {
		m := strings.Split(v, ":")
		if err := checkArray(m, name); err != nil {
			ve = append(ve, ValidationError{name, err})
			continue
		}
		switch m[0] {
		case "len":
			if l, err := strconv.Atoi(m[1]); err != nil {
				ve = append(ve, ValidationError{name, err})
			} else if len(val) > l {
				ve = append(ve, ValidationError{name, ErrLen})
			}
		case "regexp":
			if reg, err := regexp.Compile(m[1]); err != nil {
				ve = append(ve, ValidationError{name, err})
			} else if !reg.Match([]byte(val)) {
				ve = append(ve, ValidationError{name, ErrRegExp})
			}
		default:
			ve = append(ve, ValidationError{name, ErrValidate})
		}
	}
	return ve
}

func validateInt(val int, tagValue string, name string) ValidationErrors {
	var ve ValidationErrors
	if val == 0 {
		return ve
	}
	s := strings.Split(tagValue, "|")
	for _, v := range s {
		m := strings.Split(v, ":")
		if err := checkArray(m, name); err != nil {
			ve = append(ve, ValidationError{name, err})
			continue
		}
		switch m[0] {
		case "min":
			if l, err := strconv.Atoi(m[1]); err != nil {
				ve = append(ve, ValidationError{name, err})
			} else if l > val {
				ve = append(ve, ValidationError{name, ErrMin})
			}
		case "max":
			if l, err := strconv.Atoi(m[1]); err != nil {
				ve = append(ve, ValidationError{name, err})
			} else if l < val {
				ve = append(ve, ValidationError{name, ErrMax})
			}
		case "in":
			if ok, err := containsArray(strings.Split(m[1], ","), val); err != nil {
				ve = append(ve, ValidationError{name, err})
			} else if !ok {
				ve = append(ve, ValidationError{name, ErrIn})
			}
		default:
			ve = append(ve, ValidationError{name, ErrValidate})
		}
	}
	return ve
}

func validateArrInt(val []int, tagValue string, name string) ValidationErrors {
	var ve ValidationErrors
	for _, v := range val {
		if err := validateInt(v, tagValue, name); len(err) > 0 {
			ve = append(ve, err...)
		}
	}
	return ve
}

func validateArrString(val []string, tagValue string, name string) ValidationErrors {
	var ve ValidationErrors
	for _, v := range val {
		if err := validateString(v, tagValue, name); len(err) > 0 {
			ve = append(ve, err...)
		}
	}
	return ve
}
