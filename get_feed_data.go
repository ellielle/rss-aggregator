package main

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
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
func updateFeedData(cfg *apiConfig) error {
	feedList, err := cfg.getFeedsToUpdate()
	if err != nil {
		return errors.New("Error retrieving feed data")
	}

	for _, feed := range feedList {
		err := cfg.fetchFeedData(feed.Url)
		if err != nil {
			// TODO: better error message
			return errors.New("Error parsing feed data")
		}
	}

	return nil
}

// Get a list of the next <LIMIT> feeds that need to be updated
func (cfg *apiConfig) getFeedsToUpdate() ([]Feed, error) {
	const LIMIT = 10
	feeds := []Feed{}
	db_feed, err := cfg.DB.GetNextFeedsToFetch(context.Background(), LIMIT)
	if err != nil {
		return []Feed{}, errors.New(err.Error())
	}

	// Convert the feeds from the database API to a JSON-safe API
	for _, f := range db_feed {
		feed := DatabaseFeedToFeed(f)
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

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

	// TODO: update RSS URL's updated_at, and last_fetched_at fields

	return nil
}
