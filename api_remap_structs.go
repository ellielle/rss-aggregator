package main

import (
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

// table: users
type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func DatabaseUserToUser(user database.User) User {
	return User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

// table: feeds
type Feed struct {
	Id            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserId        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func DatabaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		Id:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserId:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
}

// table: feeds_follows
type Feeds_Follow struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

func DatabaseFeedFollowsToFeedFollows(ff database.FeedsFollow) Feeds_Follow {
	return Feeds_Follow{
		Id:        ff.ID,
		CreatedAt: ff.CreatedAt,
		UpdatedAt: ff.UpdatedAt,
		UserId:    ff.UserID,
		FeedId:    ff.FeedID,
	}
}
