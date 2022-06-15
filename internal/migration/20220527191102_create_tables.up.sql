CREATE TABLE IF NOT EXISTS
  feeds
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , title TEXT NOT NULL
  , description TEXT NULL CHECK(description IS NULL or length(description) > 0)
  , xml_url TEXT NOT NULL CHECK(length(xml_url) > 0)
  , html_url TEXT NULL CHECK(html_url IS NULL or length(html_url) > 0)
  , subscription_time TIMESTAMP NOT NULL DEFAULT (DATETIME('now'))
  , UNIQUE(xml_url)
);

CREATE TABLE IF NOT EXISTS
  feed_categories
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , name TEXT NOT NULL CHECK(length(name) > 0)
  , UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS
  entries
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , feed_id INTEGER NOT NULL
  , external_id TEXT NOT NULL
  , url TEXT NULL CHECK(url IS NULL or length(url) > 0)
  , title TEXT NOT NULL CHECK(title IS NULL or length(title) > 0)
  , summary TEXT NULL CHECK(summary IS NULL or length(summary) > 0)
  , content TEXT NULL CHECK(content IS NULL or length(content) > 0)
  , authors JSON NOT NULL DEFAULT '[]'
  , categories JSON NOT NULL DEFAULT '[]'
  , publication_time TIMESTAMP NULL
  , update_time TIMESTAMP NULL
  , is_read BOOLEAN NOT NULL DEFAULT false
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS entries_feed_id ON entries(feed_id);

CREATE TABLE IF NOT EXISTS
  feeds_x_feed_categories
  ( feed_id INTEGER NOT NULL
  , feed_category_id INTEGER NOT NULL
  , PRIMARY KEY (feed_id, feed_category_id)
);
