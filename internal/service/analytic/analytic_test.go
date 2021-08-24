package analytic

import (
	"fmt"
	"path/filepath"
	"testing"

	"analytics/internal/bucket"
	"analytics/internal/file_utils"
	"analytics/internal/model"
)

func BenchmarkSelect(b *testing.B) {
	// init
	path := "../../../data"
	events := make([]model.Event, 0)
	err := file_utils.ReadModel(&events, filepath.Join(path, "events.csv"))
	if err != nil {
		b.Error(err)
	}

	actors := make([]model.Actor, 0)
	err = file_utils.ReadModel(&actors, filepath.Join(path, "actors.csv"))
	if err != nil {
		b.Error(err)
	}
	repos := make([]model.Repo, 0)
	err = file_utils.ReadModel(&repos, filepath.Join(path, "repos.csv"))
	if err != nil {
		b.Error(err)
	}

	eventBucket := bucket.NewEventBucket(events)
	a := NewAnalytic(eventBucket, bucket.NewActorBucket(actors), bucket.NewRepoBucket(repos))

	b.Run("Top10UserByAmountOfPRAndCommits no index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a.Top10UserByAmountOfPRAndCommits()
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("Top10ReposByAmountOfCommits no index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a.Top10ReposByAmountOfCommits()
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("Top10ReposByAmountOfWatchEvents no index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a.Top10ReposByAmountOfWatchEvents()
			if err != nil {
				b.Error(err)
			}
		}
	})
	fmt.Println("========================================================================")

	err = eventBucket.BuildIndex("Type")
	if err != nil {
		b.Error(err)
	}
	a1 := NewAnalytic(eventBucket, bucket.NewActorBucket(actors), bucket.NewRepoBucket(repos))

	b.Run("Top10UserByAmountOfPRAndCommits index 'Type'", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a1.Top10UserByAmountOfPRAndCommits()
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("Top10ReposByAmountOfCommits index 'Type'", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a1.Top10ReposByAmountOfCommits()
			if err != nil {
				b.Error(err)
			}
		}
	})

	b.Run("Top10ReposByAmountOfWatchEvents index 'Type'", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = a1.Top10ReposByAmountOfWatchEvents()
			if err != nil {
				b.Error(err)
			}
		}
	})

	fmt.Println("========================================================================")
}
