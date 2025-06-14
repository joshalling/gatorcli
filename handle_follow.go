package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshalling/gatorcli/internal/database"
)

func handleFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the login command expects one argument: url")
	}

	user, err := s.db.GetUser(context.Background(), s.c.UserName)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed followed: %s\n", feed_follow.FeedName)
	fmt.Printf("User: %s\n", feed_follow.UserName)

	return nil
}

func handleFeedFollows(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.c.UserName)
	if err != nil {
		return err
	}

	feed_follows, err := s.db.GetFeedFollows(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed_follow := range feed_follows {
		fmt.Printf("%s\n", feed_follow.FeedName)
	}
	return nil
}
