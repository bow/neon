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
  -- subscription_time is when the feed was added into the database.
  , subscription_time TIMESTAMP NOT NULL DEFAULT (DATETIME('now'))
  -- update_time is when the feed was last updated; may be derived from update_time of the entries.
  , update_time TIMESTAMP NULL
  -- feeds must be unique by its feed URL.
  , UNIQUE(feed_url)
);

CREATE TABLE IF NOT EXISTS
  -- feed_categories lists all user-defined feed categories.
  feed_categories
  -- id is the internal ID of the feed category.
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  -- name is the string value of the category.
  , name TEXT NOT NULL CHECK(length(name) > 0)
  -- feed category must be unique by value, which is case-sensitive.
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
  -- categories contains categories linked to the entry.
  , categories JSON NOT NULL DEFAULT '[]'
  -- publication_time is when the entry was published.
  , publication_time TIMESTAMP NULL
  -- update_time is when the entry was last updated; may be the same as publication_time.
  , update_time TIMESTAMP NULL
  -- is_read indicates whether the entry has been read or not.
  , is_read BOOLEAN NOT NULL DEFAULT false
  -- entries are unique by its external ID, for a specific feed.
  , UNIQUE(feed_id, external_id)
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS entries_feed_id ON entries(feed_id);
CREATE INDEX IF NOT EXISTS entries_external_id ON entries(external_id);

CREATE TABLE IF NOT EXISTS
  -- feeds_x_feed_categories is a many-to-many table which associates feeds and feed categories.
  feeds_x_feed_categories
  -- feed_id is the database ID of the linked feed.
  ( feed_id INTEGER NOT NULL
  -- feed_category_id is the database ID of the linked feed category.
  , feed_category_id INTEGER NOT NULL
  , PRIMARY KEY (feed_id, feed_category_id)
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
  , FOREIGN KEY(feed_category_id) REFERENCES feed_categories(id) ON DELETE CASCADE
);
