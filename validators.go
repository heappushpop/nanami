package main

import (
	"fmt"
	"net/http"
)

type Field struct {
	Name   string
	Value  string
	Errors []string
}

func NewField(r *http.Request, name string, validators []Validator) *Field {
	field := &Field{
		Name:  name,
		Value: r.PostFormValue(name),
	}

	for _, validator := range validators {
		if !validator.Valid(field.Value) {
			field.Errors = append(field.Errors, validator.Message())
		}
	}

	return field
}

func (f *Field) Valid() bool {
	return len(f.Errors) == 0
}

type MaxLength struct {
	message string
	max     int
}

func NewMaxLength(max int) *MaxLength {
	return &MaxLength{
		message: fmt.Sprintf("Field length is too long, should be at most %d", max),
		max:     max,
	}
}

func (v *MaxLength) Message() string {
	return v.message
}

func (v *MaxLength) Valid(value string) bool {
	if len(value) > v.max {
		return false
	}

	return true
}

type MinLength struct {
	message string
	min     int
}

func NewMinLength(min int) *MinLength {
	return &MinLength{
		message: fmt.Sprintf("Field length is too short, should be at least %d", min),
		min:     min,
	}
}

func (v *MinLength) Message() string {
	return v.message
}

func (v *MinLength) Valid(value string) bool {
	if len(value) < v.min {
		return false
	}

	return true
}

type Validator interface {
	Message() string
	Valid(value string) bool
}
