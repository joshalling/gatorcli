# Gator CLI

Gator CLI is a command-line tool for managing users, feeds, and posts, built with Go and PostgreSQL.

## Prerequisites

- **Go** (version 1.18 or newer recommended)
- **PostgreSQL** (running locally or accessible remotely)

## Installation

Install the Gator CLI using `go install`:

```sh
go install github.com/joshalling/gatorcli@latest
```

This will place the `gatorcli` binary in your `$GOPATH/bin` or `$HOME/go/bin` directory.

## Configuration

Before running the CLI, create a configuration file in your home directory named `.gatorconfig.json`. Example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/yourdb?sslmode=disable",
  "current_user_name": ""
}
```

Replace the `db_url` with your actual PostgreSQL connection string.

## Database Setup

Run the database migrations to set up the schema:

```sh
go run github.com/pressly/goose/v3/cmd/goose@latest postgres "<your-db-url>" up
```

## Usage

Run the CLI from your terminal:

```sh
gatorcli <command> [args...]
```

### Example Commands

- `register <username>` — Register a new user.
- `login <username>` — Log in as an existing user.
- `addfeed <feed-url>` — Add a new feed (must be logged in).
- `feeds` — List all feeds.
- `browse [limit]` — Browse recent posts for the current user (default limit is 2).
- `users` — List all users.
- `reset` — Reset the users table.
- `agg` — Aggregate and scrape feeds for new posts.
- `follow <feed-id>` — Follow a feed (must be logged in).
- `unfollow <feed-id>` — Unfollow a feed (must be logged in).
- `following` — List feeds the current user is following.
