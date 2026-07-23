-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,

  -- 1. Ensure the PAIR is unique
  CONSTRAINT unique_user_feed_pair UNIQUE (user_id, feed_id),

  -- 2. Foreign Key to Users (Cascade Delete)
  CONSTRAINT fk_feed_follows_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

  -- 3. Foreign Key to Feeds (Cascade Delete)
  CONSTRAINT fk_feed_follows_feed
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;   
