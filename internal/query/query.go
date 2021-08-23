package query

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cast"
)

type Storage interface {
	Name() string
	Get(fields []string, match map[string][]interface{}) ([]map[string]interface{}, error)
}

type Query interface {
	Select(fields ...string) Query
	Count(field string) Query
	WhereIn(field string, in ...interface{}) Query
	GroupBy(field string) Query
	OrderBy(field string) Query
	Desc() Query
	Limit(int) Query
	LeftJoin(bucket Storage, left, right string) Query
	All() ([]map[string]interface{}, error)
}

type query struct {
	bucket     Storage
	fields     []string
	limit      int
	match      map[string][]interface{}
	groupBy    string
	orderBy    string
	desc       bool
	countField string
	leftJoinFn []func(row map[string]interface{}) (map[string]interface{}, error)
}

func NewQuery(bucket Storage) Query {
	return &query{bucket: bucket}
}

func (s *query) Select(fields ...string) Query {
	s.fields = fields
	return s
}

func (s *query) Count(field string) Query {
	if len(strings.TrimSpace(field)) > 0 {
		s.countField = field
	}
	return s
}

func (s *query) WhereIn(field string, in ...interface{}) Query {
	s.match = map[string][]interface{}{field: in}
	return s
}

func (s *query) GroupBy(field string) Query {
	s.groupBy = field
	return s
}

func (s *query) OrderBy(field string) Query {
	s.orderBy = field
	return s
}

func (s *query) Desc() Query {
	s.desc = true
	return s
}

func (s *query) Limit(limit int) Query {
	s.limit = limit
	return s
}

func (s *query) All() ([]map[string]interface{}, error) {
	res, err := s.bucket.Get(s.fields, s.match)
	if err != nil {
		return nil, err
	}
	if s.groupBy != "" && s.countField != "" {
		res = groupBy(s.groupBy, s.countField, res)
	} else if s.groupBy != "" {
		res = groupBy(s.groupBy, "", res)
	} else if s.countField != "" {
		res = count(s.countField, res)
	}
	if s.orderBy != "" {
		field := s.orderBy
		if strings.EqualFold(s.orderBy, "count") {
			if s.countField == "" {
				return nil, fmt.Errorf("can not order by count, the count is undefined")
			}
			field = makeCount(s.countField)
		}
		res, err = orderBy(field, s.desc, res)
		if err != nil {
			return nil, err
		}
	}
	if s.limit > 0 {
		res = limit(s.limit, res)
	}
	if s.leftJoinFn != nil {
		res, err = s.leftJoin(res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (s *query) LeftJoin(bucket Storage, left, right string) Query {
	fn := func(row map[string]interface{}) (map[string]interface{}, error) {
		if val, ok := row[left]; ok {
			rows, err := bucket.Get([]string{"*"}, map[string][]interface{}{right: {val}})
			if err != nil {
				return nil, err
			}
			if len(rows) > 0 {
				bName := bucket.Name()
				for k, v := range rows[0] {
					row[bName+"."+k] = v
				}
			}
		}
		return row, nil
	}
	if s.leftJoinFn == nil {
		s.leftJoinFn = make([]func(row map[string]interface{}) (map[string]interface{}, error), 0)
	}
	s.leftJoinFn = append(s.leftJoinFn, fn)
	return s
}

func (s *query) leftJoin(data []map[string]interface{}) ([]map[string]interface{}, error) {
	for k := range data {
		for _, f := range s.leftJoinFn {
			row, err := f(data[k])
			if err != nil {
				return nil, err
			}
			data[k] = row
		}
	}
	return data, nil
}

func count(field string, data []map[string]interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 1)
	if len(data) > 0 {
		res[0] = data[0]
		res[0][makeCount(field)] = len(data)
	}
	return res
}

func groupBy(groupByField, countField string, data []map[string]interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(data))
	deduplicate := make(map[string]int)
	for k := range data {
		if val, ok := data[k][groupByField]; ok {
			key := makeKey(groupByField, cast.ToString(val))
			if _, ok := deduplicate[key]; !ok {
				deduplicate[key] = 1
				res = append(res, data[k])
			} else {
				deduplicate[key]++
			}
		}
	}
	if countField != "" {
		for k := range res {
			if val, ok := res[k][groupByField]; ok {
				key := makeKey(groupByField, cast.ToString(val))
				if count, ok := deduplicate[key]; ok {
					res[k][makeCount(countField)] = count
				}
			}
		}
	}
	return res
}

func orderBy(field string, desc bool, data []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(data) == 0 {
		return data, nil
	}
	if _, ok := data[0][field]; !ok {
		return data, nil
	}
	switch data[0][field].(type) {
	case int, int64, int32, int16, int8:
		if desc {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToInt(data[i][field]) > cast.ToInt(data[j][field])
			})
		} else {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToInt(data[i][field]) < cast.ToInt(data[j][field])
			})
		}
	case float64, float32:
		if desc {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToFloat64(data[i][field]) > cast.ToFloat64(data[j][field])
			})
		} else {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToFloat64(data[i][field]) < cast.ToFloat64(data[j][field])
			})
		}
	case string:
		if desc {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToString(data[i][field]) > cast.ToString(data[j][field])
			})
		} else {
			sort.SliceStable(data, func(i, j int) bool {
				return cast.ToString(data[i][field]) < cast.ToString(data[j][field])
			})
		}
	default:
		return nil, fmt.Errorf("unable order by %s, unsupported field type %T to int", field, data[0][field])
	}

	return data, nil
}

func limit(limit int, data []map[string]interface{}) []map[string]interface{} {
	if len(data) > limit {
		return data[:limit]
	}
	return data
}

func makeKey(field, value string) string {
	return field + ":" + value
}

func makeCount(field string) string {
	return "count(" + field + ")"
}
