CREATE TABLE IF NOT EXISTS
  feeds
  ( id INTEGER PRIMARY KEY AUTOINCREMENT
  , ext_id TEXT NOT NULL
  , title TEXT NOT NULL
  , subtitle TEXT NULL
  , authors JSON NOT NULL DEFAULT '[]'
  , categories JSON NOT NULL DEFAULT '[]'
  , update_time TIMESTAMP NOT NULL
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
