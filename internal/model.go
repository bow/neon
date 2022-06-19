package internal

import "github.com/mmcdole/gofeed"

type Feed struct {
	gofeed.Feed
	DBID DBID
}
