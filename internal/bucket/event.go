package bucket

import (
	"fmt"

	"github.com/spf13/cast"

	"analytics/internal/model"
)

type EventBucket struct {
	data      []model.Event
	indexList map[string][]*model.Event
}

func NewEventBucket(data []model.Event) *EventBucket {
	return &EventBucket{
		data: data,
	}
}

func (s *EventBucket) Name() string {
	return "Event"
}

func (s *EventBucket) Get(fields []string, match map[string][]interface{}) ([]map[string]interface{}, error) {
	if match != nil && s.isFieldsInIndex(match) {
		return s.getIndexed(fields, match)
	} else {
		return s.fullScan(fields, match)
	}
}

// fullScan represents the full scan method
func (s *EventBucket) fullScan(fields []string, matcher map[string][]interface{}) ([]map[string]interface{}, error) {
	rows := make([]map[string]interface{}, 0)
	for _, v := range s.data {
		rowMatched := true
		row := make(map[string]interface{})
		for _, f := range fields {
			val, err := v.GetValue(f)
			if err != nil {
				return nil, err
			}
			if matcher == nil || !s.isFieldsInMather(f, matcher) {
				row[f] = val
				continue
			}
			if isValueIn(val, matcher[f]) {
				row[f] = val
			} else {
				rowMatched = false
				break
			}
		}
		if rowMatched {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

// getIndexed represents the method that retrieve the indexed data and perform the deduplication
func (s *EventBucket) getIndexed(fields []string, match map[string][]interface{}) ([]map[string]interface{}, error) {
	deduplicate := make(map[int]struct{})
	rows := make([]map[string]interface{}, 0)
	indexKeys := matcherToIndexKeys(match)
	for _, key := range indexKeys {
		if data, ok := s.indexList[key]; ok {
			for _, v := range data {
				if _, ok := deduplicate[v.ID]; ok {
					continue
				} else {
					deduplicate[v.ID] = struct{}{}
				}
				row := make(map[string]interface{})
				for _, f := range fields {
					val, err := v.GetValue(f)
					if err != nil {
						return nil, err
					}
					row[f] = val
				}
				rows = append(rows, row)
			}
		} else {
			return nil, fmt.Errorf("key " + key + " doens't exists in a index")
		}
	}
	return rows, nil
}

func isValueIn(val interface{}, list []interface{}) bool {
	for k := range list {
		wVal := cast.ToString(list[k])
		iVal := cast.ToString(val)
		if iVal != "" && iVal == wVal {
			return true
		}
	}
	return false
}

func (s *EventBucket) isFieldsInIndex(match map[string][]interface{}) (allInIndex bool) {
	allInIndex = true
	for k := range match {
		for _, val := range match[k] {
			if _, ok := s.indexList[indexKey(k, cast.ToString(val))]; !ok {
				allInIndex = false
			}
		}
	}
	return
}

func (s *EventBucket) isFieldsInMather(field string, match map[string][]interface{}) bool {
	if _, ok := match[field]; ok {
		return true
	}
	return false
}

func (s *EventBucket) BuildIndex(fields ...string) (err error) {
	if s.indexList == nil {
		s.indexList = make(map[string][]*model.Event)
	}
	for k := range s.data {
		for _, field := range fields {
			fieldVal, err := s.data[k].GetStrValue(field)
			if err != nil {
				return err
			}
			key := indexKey(field, fieldVal)
			if _, ok := s.indexList[key]; ok {
				s.indexList[key] = append(s.indexList[key], &s.data[k])
			} else {
				s.indexList[key] = []*model.Event{&s.data[k]}
			}
		}
	}
	return nil
}

func indexKey(field, value string) string {
	return field + ":" + value
}

func matcherToIndexKeys(match map[string][]interface{}) (keys []string) {
	for k := range match {
		for _, val := range match[k] {
			keys = append(keys, indexKey(k, cast.ToString(val)))
		}
	}
	return
}
