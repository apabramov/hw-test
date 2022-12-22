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
	ErrStruct   = errors.New("error struct")
)

func (v ValidationErrors) Error() string {
	var res strings.Builder
	for _, k := range v {
		res.WriteString(fmt.Sprintf("%s: %v ", k.Field, k.Err))
	}
	return res.String()
}

func Validate(v interface{}) error {
	if reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return ErrStruct
	}
	t := reflect.ValueOf(v).Elem()
	typ := t.Type()
	var ve ValidationErrors
	for i, f := range reflect.VisibleFields(typ) {
		if !f.IsExported() {
			continue
		}
		tagValue, ok := typ.Field(i).Tag.Lookup("validate")
		if !ok {
			continue
		}
		name := typ.Field(i).Name
		var err error
		switch val := t.Field(i).Interface().(type) {
		case string:
			err = validateString(val, tagValue, name)
		case []string:
			err = validateArrString(val, tagValue, name)
		case int:
			err = validateInt(val, tagValue, name)
		case []int:
			err = validateArrInt(val, tagValue, name)
		}
		var v ValidationErrors
		if !errors.As(err, &v) {
			return err
		}
		ve = append(ve, v...)
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

func validateString(val string, tagValue string, name string) error {
	var ve ValidationErrors
	if val == "" {
		return nil
	}
	s := strings.Split(tagValue, "|")
	for _, v := range s {
		m := strings.Split(v, ":")
		if err := checkArray(m, name); err != nil {
			return err
		}
		switch m[0] {
		case "len":
			if l, err := strconv.Atoi(m[1]); err != nil {
				return err
			} else if len(val) > l {
				ve = append(ve, ValidationError{name, ErrLen})
			}
		case "regexp":
			if reg, err := regexp.Compile(m[1]); err != nil {
				return err
			} else if !reg.Match([]byte(val)) {
				ve = append(ve, ValidationError{name, ErrRegExp})
			}
		default:
			ve = append(ve, ValidationError{name, ErrValidate})
		}
	}
	return ve
}

func validateInt(val int, tagValue string, name string) error {
	var ve ValidationErrors
	if val == 0 {
		return nil
	}
	s := strings.Split(tagValue, "|")
	for _, v := range s {
		m := strings.Split(v, ":")
		if err := checkArray(m, name); err != nil {
			return err
		}
		switch m[0] {
		case "min":
			if l, err := strconv.Atoi(m[1]); err != nil {
				return err
			} else if l > val {
				ve = append(ve, ValidationError{name, ErrMin})
			}
		case "max":
			if l, err := strconv.Atoi(m[1]); err != nil {
				return err
			} else if l < val {
				ve = append(ve, ValidationError{name, ErrMax})
			}
		case "in":
			if ok, err := containsArray(strings.Split(m[1], ","), val); err != nil {
				return err
			} else if !ok {
				ve = append(ve, ValidationError{name, ErrIn})
			}
		default:
			ve = append(ve, ValidationError{name, ErrValidate})
		}
	}
	return ve
}

func validateArrInt(val []int, tagValue string, name string) error {
	var ve ValidationErrors
	for _, v := range val {
		if err := validateInt(v, tagValue, name); err != nil {
			var v ValidationErrors
			if !errors.As(err, &v) {
				return err
			}
			ve = append(ve, v...)
		}
	}
	return ve
}

func validateArrString(val []string, tagValue string, name string) error {
	var ve ValidationErrors
	for _, v := range val {
		if err := validateString(v, tagValue, name); err != nil {
			var v ValidationErrors
			if !errors.As(err, &v) {
				return err
			}
			ve = append(ve, v...)
		}
	}
	return ve
}
