package model

import (
	"fmt"
	"strconv"
)

type Actor struct {
	ID       int    `csv:"id"`
	UserName string `csv:"username"`
}

func (s *Actor) GetStrValue(fieldName string) (string, error) {
	switch fieldName {
	case "ID":
		return strconv.Itoa(s.ID), nil
	case "UserName":
		return s.UserName, nil
	default:
		return "", fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}

func (s *Actor) GetValue(fieldName string) (interface{}, error) {
	switch fieldName {
	case "ID":
		return s.ID, nil
	case "UserName":
		return s.UserName, nil
	default:
		return nil, fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}
