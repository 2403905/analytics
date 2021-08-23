package bucket

import (
	"fmt"
	"strconv"

	"github.com/spf13/cast"

	"analytics/internal/model"
)

type ActorBucket struct {
	data    []model.Actor
	pkIndex map[string]*model.Actor
}

func NewActorBucket(data []model.Actor) *ActorBucket {
	pkIndex := make(map[string]*model.Actor)
	for k := range data {
		pkIndex[strconv.Itoa(data[k].ID)] = &data[k]
	}
	return &ActorBucket{
		data:    data,
		pkIndex: pkIndex,
	}
}

func (s *ActorBucket) Name() string {
	return "Actor"
}

func (s *ActorBucket) Get(fields []string, match map[string][]interface{}) ([]map[string]interface{}, error) {
	if m, ok := match["ID"]; ok && len(m) == 1 {
		if row, ok := s.pkIndex[cast.ToString(m[0])]; ok {
			return []map[string]interface{}{{"UserName": row.UserName}}, nil
		}
	}
	return nil, fmt.Errorf("the Actor bucket suppots math by a sinle ID")
}
