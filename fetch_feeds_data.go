package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
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
	feedList, err := cfg.getFeedsToUpdate()
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
				err = cfg.fetchFeedData(URL)
				if err != nil {
					return errors.New("Bad XML feed")
				}

				// Update last_fetched_at and updated_at fields of the Feed
				cfg.updateFeedMetadata(feed.Id)
				return err
			}(feed.Url)
		}
		// Wait for all feed fetches to complete
		wg.Wait()
	}
}

// Get a list of the next <LIMIT> feeds that need to be updated
func (cfg *apiConfig) getFeedsToUpdate() ([]GetNextFeedsToFetchRow, error) {
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

// Run in a goroutine so it can be concurrently processed
// while handling http requests
//

// Update feed data for a single feed URL
func (cfg *apiConfig) fetchFeedData(URL string) error {
	// Create a request to the RSS feed's URL
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return errors.New("Invalid feed URL")
	}

	// Carry out the http request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Invalid feed URL")
	}
	defer resp.Body.Close()

	// Get body as []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Invalid RSS feed")
	}

	// Shove that data into that struct
	feedData := RSSFeed{}
	err = xml.Unmarshal(body, &feedData)
	if err != nil {
		return errors.New("Invalid RSS feed")
	}

	// TODO: this can be removed when done, it's just a log
	log.Printf("feedData: %v", feedData.Channel.Item[0].Title)

	return nil
}

func (cfg *apiConfig) updateFeedMetadata(feedID uuid.UUID) {
	fetched := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: feedID, LastFetchedAt: fetched})
}
