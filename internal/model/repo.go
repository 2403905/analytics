package model

import (
	"fmt"
	"strconv"
)

type Repo struct {
	ID   int    `csv:"id"`
	Name string `csv:"name"`
}

func (s *Repo) GetStrValue(fieldName string) (string, error) {
	switch fieldName {
	case "ID":
		return strconv.Itoa(s.ID), nil
	case "Name":
		return s.Name, nil
	default:
		return "", fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}

func (s *Repo) GetValue(fieldName string) (interface{}, error) {
	switch fieldName {
	case "ID":
		return s.ID, nil
	case "Name":
		return s.Name, nil
	default:
		return nil, fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}
