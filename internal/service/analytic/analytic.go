package analytic

import (
	"analytics/internal/query"
)

type Analytic struct {
	eventBucket query.Storage
	actorBucket query.Storage
	repoBucket  query.Storage
}

func NewAnalytic(eventBucket, actorBucket, repoBucket query.Storage) Analytic {
	return Analytic{
		eventBucket: eventBucket,
		actorBucket: actorBucket,
		repoBucket:  repoBucket,
	}
}

func (s *Analytic) Top10UserByAmountOfPRAndCommits() ([]map[string]interface{}, error) {
	return query.NewQuery(s.eventBucket).Select("ID", "Type", "ActorID").
		WhereIn("Type", "PullRequestEvent", "PushEvent").
		Count("Type").
		GroupBy("ActorID").
		OrderBy("Count").
		Desc().
		Limit(10).
		LeftJoin(s.actorBucket, "ActorID", "ID").
		All()
}

func (s Analytic) Top10ReposByAmountOfCommits() ([]map[string]interface{}, error) {
	return query.NewQuery(s.eventBucket).Select("ID", "Type", "RepoID").
		WhereIn("Type", "PushEvent").
		Count("Type").
		GroupBy("RepoID").
		OrderBy("Count").
		Desc().
		LeftJoin(s.repoBucket, "RepoID", "ID").
		Limit(10).All()
}

func (s Analytic) Top10ReposByAmountOfWatchEvents() ([]map[string]interface{}, error) {
	return query.NewQuery(s.eventBucket).Select("ID", "Type", "RepoID").
		WhereIn("Type", "WatchEvent").
		Count("Type").
		GroupBy("RepoID").
		OrderBy("Count").
		Desc().
		LeftJoin(s.repoBucket, "RepoID", "ID").
		Limit(10).All()
}
