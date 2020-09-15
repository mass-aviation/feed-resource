package types

import (
	"time"
)

type Version struct {
        PubDate string `json:"pubDate"`
        PubDateParsed *time.Time `json:"-"`
}

type InputReq struct {
        Source struct {
                Url string `json:"url"`
        } `json:"source"`
        Version *Version `json:"version"`
}

type Output struct {
        Version Version `json:"version"`
        Metadata []map[string]string `json:"metadata"`
}
