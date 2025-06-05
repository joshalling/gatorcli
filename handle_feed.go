package main

import (
	"context"
	"fmt"

	"github.com/joshalling/gatorcli/internal/api"
)

func handleAgg(s *state, cmd command) error {
	feed, err := api.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}
