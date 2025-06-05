package api

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(ctx context.Context, url string) (RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return RSSFeed{}, err
	}

	request.Header.Set("User-Agent", "gatorcli")
	request.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return RSSFeed{}, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return RSSFeed{}, err
	}
	return feed, nil
}
