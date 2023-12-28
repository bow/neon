CREATE TABLE IF NOT EXISTS
  -- feeds contains all subscribed feeds.
  feeds
  -- id is the internal database ID of the feed.
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  -- title is the title of the feed; may be user-defined.
  , title TEXT NOT NULL
  -- description is the feed description; may be user-defined.
  , description TEXT NULL CHECK(description IS NULL or length(description) > 0)
  -- feed_url is the URL of the of the feed.
  , feed_url TEXT NOT NULL CHECK(length(feed_url) > 0)
  -- site_url is the URL of the website linked to the feed.
  , site_url TEXT NULL CHECK(site_url IS NULL or length(site_url) > 0)
  -- is_starred indicates whether the entry has been starred or not.
  , is_starred BOOLEAN NOT NULL DEFAULT false
  -- sub_time is when the feed was added into the database.
  , sub_time TIMESTAMP NOT NULL DEFAULT (DATETIME('now'))
  -- update_time is when the feed was last updated; may be derived from update_time of the entries.
  , update_time TIMESTAMP NULL
  -- last_pull_time is when the feed was last pulled.
  , last_pull_time TIMESTAMP NOT NULL
  -- feeds must be unique by its feed URL.
  , UNIQUE(feed_url)
  );

CREATE TABLE IF NOT EXISTS
  -- feed_tags lists all user-defined feed tags.
  feed_tags
  -- id is the internal ID of the feed tag.
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  -- name is the string value of the tag.
  , name TEXT NOT NULL CHECK(length(name) > 0)
  -- feed tag must be unique by value, which is case-sensitive.
  , UNIQUE(name)
  );

CREATE TABLE IF NOT EXISTS
  -- entries contains all entries linked to a given feed.
  entries
  -- id is the internal database ID of the entry.
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  -- feed_id is the internal database ID of the feed to which the entry belongs.
  , feed_id INTEGER NOT NULL
  -- external_id is the externally-defined ID value of the entry
  , external_id TEXT NOT NULL
  -- url is the URL to which the entry is linked.
  , url TEXT NULL CHECK(url IS NULL or length(url) > 0)
  -- title is the entry title.
  , title TEXT NOT NULL CHECK(title IS NULL or length(title) > 0)
  -- description is the entry description.
  , description TEXT NULL CHECK(description IS NULL or length(description) > 0)
  -- contents is the actual content of the entry.
  , content TEXT NULL CHECK(content IS NULL or length(content) > 0)
  -- authors lists the authors and contributors of the entry.
  , authors JSON NOT NULL DEFAULT '[]'
  -- tags contains tags linked to the entry.
  , tags JSON NOT NULL DEFAULT '[]'
  -- pub_time is when the entry was published.
  , pub_time TIMESTAMP NULL
  -- update_time is when the entry was last updated; may be the same as pub_time.
  , update_time TIMESTAMP NULL
  -- is_read indicates whether the entry has been read or not.
  , is_read BOOLEAN NOT NULL DEFAULT false
  -- is_bookmarked indicates the bookmark status of the entry.
  , is_bookmarked BOOLEAN NOT NULL DEFAULT false
  -- entries are unique by its external ID, for a specific feed.
  , UNIQUE(feed_id, external_id)
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
  );
CREATE INDEX IF NOT EXISTS entries_feed_id ON entries(feed_id);
CREATE INDEX IF NOT EXISTS entries_external_id ON entries(external_id);

CREATE TABLE IF NOT EXISTS
  -- feeds_x_feed_tags is a many-to-many table which associates feeds and feed tags.
  feeds_x_feed_tags
  -- feed_id is the database ID of the linked feed.
  ( feed_id INTEGER NOT NULL
  -- feed_tag_id is the database ID of the linked feed tag.
  , feed_tag_id INTEGER NOT NULL
  , PRIMARY KEY (feed_id, feed_tag_id)
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
  , FOREIGN KEY(feed_tag_id) REFERENCES feed_tags(id) ON DELETE CASCADE
  );
