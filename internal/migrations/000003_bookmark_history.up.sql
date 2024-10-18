BEGIN;

CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    visible BOOLEAN NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_bookmarks_name UNIQUE (name),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);


CREATE TABLE IF NOT EXISTS bookmark_blogs (
    bookmark_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_bookmark_blogs_bookmark_id_blog_id UNIQUE (bookmark_id, blog_id),
    FOREIGN KEY (bookmark_id) REFERENCES bookmarks(id) ON DELETE CASCADE,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS read_later (
    user_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_read_later_user_id_blog_id UNIQUE (user_id, blog_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS reading_history (
    user_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_reading_history_user_id_blog_id UNIQUE (user_id, blog_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);


COMMIT;