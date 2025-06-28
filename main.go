package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joshalling/gatorcli/internal/config"
	"github.com/joshalling/gatorcli/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	c  *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", c.DbUrl)
	if err != nil {
		fmt.Printf("Error connecting to database: %v", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := &state{
		db: dbQueries,
		c:  &c,
	}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleList)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleCreateFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("following", middlewareLoggedIn(handleFeedFollows))

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	err = cmds.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.c.UserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
