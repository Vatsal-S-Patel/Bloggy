BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  username VARCHAR(30) NOT NULL UNIQUE,
  email VARCHAR(70) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  bio TEXT,
  avatar TEXT,
  followers INT,
  following INT,
  joined_at TIMESTAMPTZ NOT NULL,
  last_login_at TIMESTAMPTZ NOT NULL
);


CREATE TABLE IF NOT EXISTS admins(
  id UUID PRIMARY KEY,
  username VARCHAR(30) NOT NULL UNIQUE,
  email VARCHAR(70) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  avatar TEXT,
  joined_at TIMESTAMPTZ NOT NULL,
  last_login_at TIMESTAMPTZ NOT NULL
);


CREATE TABLE IF NOT EXISTS superadmins(
  id UUID PRIMARY KEY,
  username VARCHAR(30) NOT NULL UNIQUE,
  email VARCHAR(70) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  joined_at TIMESTAMPTZ NOT NULL,
  last_login_at TIMESTAMPTZ NOT NULL
);


CREATE TABLE IF NOT EXISTS followers (
  follower_id UUID NOT NULL,
  following_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT unique_follower_following UNIQUE (follower_id, following_id)
);

CREATE INDEX IF NOT EXISTS idx_follower_id_created_at ON followers(follower_id, created_at);

CREATE INDEX IF NOT EXISTS idx_following_id_created_at ON followers(following_id, created_at);

COMMIT;