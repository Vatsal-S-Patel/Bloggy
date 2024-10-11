BEGIN;

CREATE TABLE IF NOT EXISTS blogs (
  id UUID PRIMARY KEY,
  title VARCHAR(130) NOT NULL,
  subtitle VARCHAR(170),
  content TEXT NOT NULL,
  ft_image TEXT,
  author_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_blogs_title ON blogs(title);


CREATE TABLE IF NOT EXISTS tags (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS blog_tags (
    blog_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    CONSTRAINT unique_tag_blog UNIQUE (tag_id, blog_id)
);

CREATE INDEX IF NOT EXISTS idx_blog_tags_blog_id ON blog_tags(blog_id);


COMMIT;