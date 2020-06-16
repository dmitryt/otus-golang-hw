//Code generated by cool go-validate tool; DO NOT EDIT.
package models

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

type _validateValidator func() error

func (e ValidationError) Error() string { return fmt.Sprintf("%s: %s;", e.Field, e.Err.Error()) }

var (
	ErrValidationLen      = errors.New("Incorrect value length")
	ErrValidationMin      = errors.New("Value is less than minimum")
	ErrValidationMax      = errors.New("Value is greater than maximum")
	ErrValidationContains = errors.New("Value doesn't contain in provided list")
	ErrValidationRegexp   = errors.New("value doesn't match by pattern")
)

func _helperContainsStr(val string, items string) bool {
	arr := strings.Split(items, ",")
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func _helperContainsInt(val int, items string) (bool, error) {
	arr := strings.Split(items, ",")
	for _, item := range arr {
		result, err := strconv.Atoi(item)
		if err != nil {
			return false, err
		}
		if result == val {
			return true, nil
		}
	}
	return false, nil
}

func (model User) _validateLenID() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("ID").IsZero() {
		return nil
	}

	for _, value := range []string{model.ID} {

		if len(value) != 36 {
			return ValidationError{Field: "ID", Err: ErrValidationLen}
		}

	}
	return nil
}

func (model User) _validateMinAge() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Age").IsZero() {
		return nil
	}

	for _, value := range []int{model.Age} {

		if value < 18 {
			return ValidationError{Field: "Age", Err: ErrValidationMin}
		}

	}
	return nil
}

func (model User) _validateMaxAge() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Age").IsZero() {
		return nil
	}

	for _, value := range []int{model.Age} {

		if value > 50 {
			return ValidationError{Field: "Age", Err: ErrValidationMax}
		}

	}
	return nil
}

func (model User) _validateRegexpEmail() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Email").IsZero() {
		return nil
	}

	for _, value := range []string{model.Email} {

		re := regexp.MustCompile("^\\w+@\\w+\\.\\w+$")
		if !re.Match([]byte(value)) {
			return ValidationError{Field: "Email", Err: ErrValidationRegexp}
		}

	}
	return nil
}

func (model User) _validateInStringRole() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Role").IsZero() {
		return nil
	}

	for _, value := range []UserRole{model.Role} {

		if !_helperContainsStr(string(value), "admin,stuff") {
			return ValidationError{Field: "Role", Err: ErrValidationContains}
		}

	}
	return nil
}

func (model User) _validateLenPhones() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Phones").IsZero() {
		return nil
	}

	for _, value := range model.Phones {

		if len(value) != 11 {
			return ValidationError{Field: "Phones", Err: ErrValidationLen}
		}

	}
	return nil
}

func (model User) Validate() ([]ValidationError, error) {
	var vErrors []ValidationError
	validators := []_validateValidator{

		model._validateLenID,

		model._validateMinAge,

		model._validateMaxAge,

		model._validateRegexpEmail,

		model._validateInStringRole,

		model._validateLenPhones,
	}
	for _, fn := range validators {
		err := fn()
		if err == nil {
			continue
		}
		if vError, ok := err.(ValidationError); ok {
			vErrors = append(vErrors, vError)
		} else {
			return []ValidationError{}, err
		}
	}
	return vErrors, nil
}

func (model App) _validateLenVersion() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Version").IsZero() {
		return nil
	}

	for _, value := range []string{model.Version} {

		if len(value) != 5 {
			return ValidationError{Field: "Version", Err: ErrValidationLen}
		}

	}
	return nil
}

func (model App) Validate() ([]ValidationError, error) {
	var vErrors []ValidationError
	validators := []_validateValidator{

		model._validateLenVersion,
	}
	for _, fn := range validators {
		err := fn()
		if err == nil {
			continue
		}
		if vError, ok := err.(ValidationError); ok {
			vErrors = append(vErrors, vError)
		} else {
			return []ValidationError{}, err
		}
	}
	return vErrors, nil
}

func (model Response) _validateInIntCode() error {
	rv := reflect.ValueOf(&model)
	//field is not set
	if rv.Elem().FieldByName("Code").IsZero() {
		return nil
	}

	for _, value := range []int{model.Code} {

		result, err := _helperContainsInt(value, "200,404,500")
		if err != nil {
			return err
		}
		if !result {
			return ValidationError{Field: "Code", Err: ErrValidationContains}
		}

	}
	return nil
}

func (model Response) Validate() ([]ValidationError, error) {
	var vErrors []ValidationError
	validators := []_validateValidator{

		model._validateInIntCode,
	}
	for _, fn := range validators {
		err := fn()
		if err == nil {
			continue
		}
		if vError, ok := err.(ValidationError); ok {
			vErrors = append(vErrors, vError)
		} else {
			return []ValidationError{}, err
		}
	}
	return vErrors, nil
}