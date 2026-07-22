package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (f *RSSFeed) UnescapeHTML() {
	// Process Channel fields
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	// Process Item fields
	for i := range f.Channel.Item {
		f.Channel.Item[i].Title = html.UnescapeString(f.Channel.Item[i].Title)
		f.Channel.Item[i].Description = html.UnescapeString(f.Channel.Item[i].Description)
	}
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		errMsg := fmt.Errorf("failed to create request %s: %v", feedURL, err)
		return nil, errMsg
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Errorf("failed http request for %s: %v", feedURL, err)
		return nil, errMsg
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Errorf("failed to read resp bytes: %v", err)
		return nil, errMsg
	}
	defer resp.Body.Close()

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		errMsg := fmt.Errorf("failed to marshal xml: %v", err)
		return nil, errMsg
	}
	feed.UnescapeHTML()

	return &feed, nil
}
