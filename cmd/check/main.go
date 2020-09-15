package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/mass-aviation/feed-resource/internal/types"
)

func main() {
	var err error
	var inputReq types.InputReq
	if err := json.NewDecoder(os.Stdin).Decode(&inputReq); err != nil {
		panic(err)
	}

	// Validate input
	if inputReq.Source.Url == "" {
		panic("feed url is not set")
	}

	lastDate := time.Unix(0, 0)
	if v := inputReq.Version; v != nil {
		lastDate, err = time.Parse(time.RFC3339, inputReq.Version.PubDate);
		if err != nil {
			panic(err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(inputReq.Source.Url, ctx)
	if err != nil {
		panic(err)
	}

	cancel()

	collected := []*gofeed.Item{}
	for _, item := range feed.Items {
		if item.PublishedParsed.After(lastDate) {
			collected = append(collected, item)
		}
	}

	// Reorder the versions
	versions := []*types.Version{}
	for i := len(collected)-1; i >= 0; i-- {
		item := collected[i]
		versions = append(versions, &types.Version{
			PubDate: item.PublishedParsed.Format(time.RFC3339),
		})
	}

	if err := json.NewEncoder(os.Stdout).Encode(versions); err != nil {
		panic(err)
	}
}
