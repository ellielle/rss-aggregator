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
type FeedsFollow struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	FeedId    uuid.UUID `json:"feed_id"`
}

func DatabaseFeedFollowsToFeedFollows(ff database.FeedsFollow) FeedsFollow {
	return FeedsFollow{
		Id:        ff.ID,
		CreatedAt: ff.CreatedAt,
		UpdatedAt: ff.UpdatedAt,
		UserId:    ff.UserID,
		FeedId:    ff.FeedID,
	}
}

type GetNextFeedsToFetchRow struct {
	Id  uuid.UUID `json:"id"`
	Url string    `json:"url"`
}

func DatabaseNextFeedsToNextFeeds(nf database.GetNextFeedsToFetchRow) GetNextFeedsToFetchRow {
	return GetNextFeedsToFetchRow{
		Id:  nf.ID,
		Url: nf.Url,
	}
}

// table: posts
type Post struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description *string   `json:"description"`
	FeedId      uuid.UUID `json:"feed_id"`
}

func DatabasePostsToPosts(post database.Post) Post {
	return Post{
		Id:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: &post.Description.String,
		FeedId:      post.FeedID,
	}
}
