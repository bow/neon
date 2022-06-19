package internal

import "github.com/mmcdole/gofeed"

type Feed struct {
	DBID  DBID
	inner gofeed.Feed
}
