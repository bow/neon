CREATE TABLE IF NOT EXISTS
  feeds
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , title TEXT NOT NULL
  , description TEXT NULL
  , creation_time TIMESTAMP NOT NULL DEFAULT (DATETIME('now'))
);

CREATE TABLE IF NOT EXISTS
  feed_categories
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , name TEXT NOT NULL
  , UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS
  entries
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , feed_id INTEGER NOT NULL
  , ext_id TEXT NOT NULL
  , title TEXT NOT NULL
  , content TEXT NULL
  , authors JSON NOT NULL DEFAULT '[]'
  , categories JSON NOT NULL DEFAULT '[]'
  , publish_time TIMESTAMP NULL
  , update_time TIMESTAMP NOT NULL
  , FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS
  feeds_x_feed_categories
  ( feed_id INTEGER NOT NULL
  , feed_category_id INTEGER NOT NULL
  , PRIMARY KEY (feed_id, feed_category_id)
);
