package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bow/iris/internal/store"
)

func newFeedListCmd() *cobra.Command {
	var (
		name = "list"
	)

	list := cobra.Command{
		Use:   name,
		Short: "List feeds",
		RunE: func(cmd *cobra.Command, args []string) error {

			str, err := storeFromCtx(cmd)
			if err != nil {
				return err
			}

			feeds, err := str.ListFeeds(cmd.Context())
			if err != nil {
				return err
			}
			for _, feed := range feeds {
				fmt.Printf("%s", fmtFeed(feed))
			}

			return nil
		},
	}

	return &list
}

func fmtFeed(feed *store.Feed) string {
	var (
		sb  strings.Builder
		cat = func(format string, a ...any) { fmt.Fprintf(&sb, format, a...) }
	)

	var upds = "?"
	upd, err := store.DeserializeTime(&feed.Updated.String)
	if err != nil {
		upd = nil
	}
	if upd != nil {
		upds = upd.Local().Format("2 January 2006 • 15:04 MST")
	}

	var nread, ntotal int
	for _, entry := range feed.Entries {
		if entry.IsRead {
			nread++
		}
		ntotal++
	}

	cat("\x1b[36m▶\x1b[0m \x1b[4m%s\x1b[0m\n", capText(feed.Title))
	cat("  ID     : %d\n", feed.DBID)
	cat("  Updated: %s\n", upds)
	cat("  Unread : %d/%d\n", ntotal-nread, ntotal)
	cat("  URL    : %s\n", capText(feed.SiteURL.String))
	cat("  Tags   : #%s\n", strings.Join(feed.Tags, " #"))
	cat("\n")

	return sb.String()
}

func capText(txt string) string {
	if len(txt) > (displayWidth - indentWidth) {
		return fmt.Sprintf("%s%s", txt[:displayWidth-indentWidth-len(ellipsis)], ellipsis)
	}
	return txt
}
