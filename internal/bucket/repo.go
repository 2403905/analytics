package bucket

import (
	"fmt"
	"strconv"

	"github.com/spf13/cast"

	"analytics/internal/model"
)

type RepoBucket struct {
	data    []model.Repo
	pkIndex map[string]*model.Repo
}

func NewRepoBucket(data []model.Repo) *RepoBucket {
	pkIndex := make(map[string]*model.Repo)
	for k := range data {
		pkIndex[strconv.Itoa(data[k].ID)] = &data[k]
	}
	return &RepoBucket{
		data:    data,
		pkIndex: pkIndex,
	}
}

func (s *RepoBucket) Name() string {
	return "Repo"
}

func (s *RepoBucket) Get(fields []string, match map[string][]interface{}) ([]map[string]interface{}, error) {
	if m, ok := match["ID"]; ok && len(m) == 1 {
		if row, ok := s.pkIndex[cast.ToString(m[0])]; ok {
			return []map[string]interface{}{{"UserName": row.Name}}, nil
		}
	}
	return nil, fmt.Errorf("the Repo bucket suppots math by a sinle ID")
}
