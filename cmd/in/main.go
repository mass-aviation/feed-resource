package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/mass-aviation/feed-resource/internal/types"
)

func main() {
	destDir := ""
	if len(os.Args) > 1 {
		destDir = os.Args[1]
	} else {
		panic("destDir not specified")
	}

	var err error
	var inputReq types.InputReq
	if err := json.NewDecoder(os.Stdin).Decode(&inputReq); err != nil {
		panic(err)
	}

	// Validate input
	if inputReq.Source.Url == "" {
		panic("feed url is not set")
	}

	wantedDate := time.Unix(0, 0)
	if v := inputReq.Version; v != nil {
		wantedDate, err = time.Parse(time.RFC3339, inputReq.Version.PubDate);
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

	var foundItem *gofeed.Item = nil
	for _, item := range feed.Items {
		if item.PublishedParsed.Equal(wantedDate) {
			foundItem = item
			break
		}
	}

	if foundItem == nil {
		panic("could not find specified item from the feed")
	}

	ensure(writeItem(destDir, "title", foundItem.Title))
	ensure(writeItem(destDir, "content", foundItem.Content))
	ensure(writeItem(destDir, "description", foundItem.Description))
	ensure(writeItem(destDir, "guid", foundItem.GUID))
	ensure(writeItem(destDir, "link", foundItem.Link))
	ensure(writeItem(destDir, "pubDate", foundItem.PublishedParsed.Format(time.RFC3339)))

	out := types.Output{
		Version: *inputReq.Version,
		Metadata: []map[string]string{
			{"name": "title", "value": foundItem.Title},
			{"name": "description", "value": foundItem.Description},
			{"name": "guid", "value": foundItem.GUID},
		},
	}

	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		panic(err)
	}
}

func writeItem(dir string, name string, content string) error {
	dest := filepath.Join(dir, name)

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return err
	}

	return nil
}

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}
