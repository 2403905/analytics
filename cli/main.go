package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"analytics/internal/bucket"
	"analytics/internal/file_utils"
	"analytics/internal/model"
	"analytics/internal/service/analytic"
)

func main() {
	path := "data"
	events := make([]model.Event, 0)
	err := file_utils.ReadModel(&events, filepath.Join(path, "events.csv"))
	if err != nil {
		log.Err(err).Send()
	}

	actors := make([]model.Actor, 0)
	err = file_utils.ReadModel(&actors, filepath.Join(path, "actors.csv"))
	if err != nil {
		log.Err(err).Send()
	}
	repos := make([]model.Repo, 0)
	err = file_utils.ReadModel(&repos, filepath.Join(path, "repos.csv"))
	if err != nil {
		log.Err(err).Send()
	}

	eventBucket := bucket.NewEventBucket(events)
	err = eventBucket.BuildIndex("Type")
	if err != nil {
		log.Err(err).Send()
	}

	a := analytic.NewAnalytic(eventBucket, bucket.NewActorBucket(actors), bucket.NewRepoBucket(repos))

	for {
		fmt.Println("Choose an option:")
		fmt.Println("[1] Top 10 active users sorted by amount of PRs created and commits pushed")
		fmt.Println("[2] Top 10 repositories sorted by amount of commits pushed")
		fmt.Println("[3] Top 10 repositories sorted by amount of watch events")
		fmt.Println("[0] Exit")
		var i int
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			log.Err(err).Send()
			continue
		}
		switch i {
		case 1:
			top10, err := a.Top10UserByAmountOfPRAndCommits()
			if err != nil {
				log.Err(err).Send()
			}
			prettyPrint(top10)
		case 2:
			top10, err := a.Top10ReposByAmountOfCommits()
			if err != nil {
				log.Err(err).Send()
			}
			prettyPrint(top10)
		case 3:
			top10, err := a.Top10ReposByAmountOfWatchEvents()
			if err != nil {
				log.Err(err).Send()
			}
			prettyPrint(top10)
		case 0:
			os.Exit(0)
		}
	}
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	} else {
		log.Err(err).Send()
	}
	return
}
