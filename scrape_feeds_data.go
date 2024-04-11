package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ellielle/rss-aggregator/internal/database"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

// Update an RSS feed's data
func updateFeedData(cfg *apiConfig, tickerInterval time.Duration) {
	feedList, err := getFeedsToUpdate(cfg)
	if err != nil {
		return
	}

	ticker := time.NewTicker(tickerInterval)
	for ; ; <-ticker.C {

		// Use a WaitGroup to update all feeds at the same time
		wg := sync.WaitGroup{}
		for _, feed := range feedList {
			// Add 1 to the tracked goroutines
			wg.Add(1)

			go func(URL string) error {
				// Decrement tracked goroutines when finished
				defer wg.Done()

				// Attempt to update RSS Feed
				feedData, err := fetchFeedData(URL)
				if err != nil {
					return errors.New("Bad XML feed")
				}

				// Update database with Posts data
				updateFeedsDatabase(cfg, feedData, feed.Id)

				// Update last_fetched_at and updated_at fields of the Feed
				updateFeedMetadata(cfg, feed.Id)
				return err
			}(feed.Url)
		}
		// Wait for all feed fetches to complete
		wg.Wait()
	}
}

// Get a list of the next <LIMIT> feeds that need to be updated
func getFeedsToUpdate(cfg *apiConfig) ([]GetNextFeedsToFetchRow, error) {
	const LIMIT = 10
	feeds, err := cfg.DB.GetNextFeedsToFetch(context.Background(), LIMIT)
	if err != nil {
		return []GetNextFeedsToFetchRow{}, errors.New(err.Error())
	}

	// Map the database API struct to a JSON-friendly struct
	mfeeds := []GetNextFeedsToFetchRow{}
	for _, feed := range feeds {
		f := DatabaseNextFeedsToNextFeeds(feed)
		mfeeds = append(mfeeds, f)
	}

	return mfeeds, nil
}

// Update feed data for a single feed URL
func fetchFeedData(URL string) (*RSSFeed, error) {
	// Create a request to the RSS feed's URL
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return &RSSFeed{}, errors.New("Invalid feed URL")
	}

	// Carry out the http request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, errors.New("Invalid feed URL")
	}
	defer resp.Body.Close()

	// Get body as []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, errors.New("Invalid RSS feed")
	}

	// Shove that data into that struct
	feedData := RSSFeed{}
	err = xml.Unmarshal(body, &feedData)
	if err != nil {
		return &RSSFeed{}, errors.New("Invalid RSS feed")
	}

	return &feedData, nil
}

func updateFeedMetadata(cfg *apiConfig, feedID uuid.UUID) {
	fetched := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: feedID, LastFetchedAt: fetched})
}

// Update database with current Feed data
func updateFeedsDatabase(cfg *apiConfig, data *RSSFeed, feedId uuid.UUID) error {
	// iterate over each post in the feed and save it

	for _, post := range data.Channel.Item {
		// generate a UUID for each post
		id := uuid.New()

		// Parse RFC1123Z time layout string into time.Time
		postCreatedTime, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			log.Printf("Error parsing time %v for post with id %v", post.PubDate, feedId)
		}

		// Deal with sql nullstring
		description := stringToSqlNullString(post.Description)

		_, err = cfg.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          id,
			CreatedAt:   postCreatedTime,
			UpdatedAt:   postCreatedTime,
			Title:       post.Title,
			Url:         post.Link,
			Description: description,
			FeedID:      feedId,
		})
		// Ignore SQL error of duplicate URL, Post is already in database
		if err != nil && !strings.Contains(err.Error(), "posts_url_key") {
			return errors.New(err.Error())
		}
	}

	return nil
}

func stringToSqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
