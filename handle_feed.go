package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joshalling/gatorcli/internal/api"
	"github.com/joshalling/gatorcli/internal/database"
	"github.com/lib/pq"
)

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("the agg command expects one arguments: time_between_reqs")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for range ticker.C {
		fmt.Println("========")
		fmt.Println("Scraping feeds...")
		fmt.Println("========")
		scrapeFeeds(s)
		fmt.Println("========")
		fmt.Println("Ending scrape...")
		fmt.Println("========")
	}
	return nil
}

func handleCreateFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("the add feed command expects two arguments: name, url")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      cmd.args[0],
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Created feed %s\n", feed.Name)
	fmt.Printf("Feed ID: %s\n", feed.ID)
	fmt.Printf("Created at: %s\n", feed.CreatedAt)
	fmt.Printf("Updated at: %s\n", feed.UpdatedAt)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("User ID: %s\n", feed.UserID)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s | %s | %s\n", feed.Name, feed.Url, feed.Name_2)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	rssFeed, err := api.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	for _, item := range rssFeed.Channel.Items {
		publishedAt, _ := parsePublishedAt(item.PubDate)
		err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: item.Title != ""},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullTime{Time: publishedAt, Valid: !publishedAt.IsZero()},
			FeedID:      feed.ID,
		})
		if err != nil && !isUniqueViolation(err) {
			fmt.Printf("Error saving post: %v\n", err)
		}
	}

	return nil
}

func handleBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		fmt.Sscanf(cmd.args[0], "%d", &limit)
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s\n%s\n%s\nPublished: %s\n---\n", post.Title.String, post.Url, post.Description.String, post.PublishedAt.Time.Format(time.RFC3339))
	}
	return nil
}

func parsePublishedAt(pubDate string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z, time.RFC1123, time.RFC822Z, time.RFC822,
		time.RFC3339, time.RFC3339Nano,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, pubDate); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time: %s", pubDate)
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
