package model

import (
	"fmt"
	"strconv"
)

type Event struct {
	ID      int    `csv:"id"`
	Type    string `csv:"type"`
	ActorID int    `csv:"actor_id"`
	RepoID  int    `csv:"repo_id"`
}

func (s *Event) GetStrValue(fieldName string) (string, error) {
	switch fieldName {
	case "ID":
		return strconv.Itoa(s.ID), nil
	case "Type":
		return s.Type, nil
	case "ActorID":
		return strconv.Itoa(s.ActorID), nil
	case "RepoID":
		return strconv.Itoa(s.RepoID), nil
	default:
		return "", fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}

func (s *Event) GetValue(fieldName string) (interface{}, error) {
	switch fieldName {
	case "ID":
		return s.ID, nil
	case "Type":
		return s.Type, nil
	case "ActorID":
		return s.ActorID, nil
	case "RepoID":
		return s.RepoID, nil
	default:
		return nil, fmt.Errorf("field " + fieldName + " doesn't exist")
	}
}
