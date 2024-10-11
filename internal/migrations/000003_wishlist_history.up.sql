BEGIN;

CREATE TABLE IF NOT EXISTS wishlists (
    user_id UUID,
    blog_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, blog_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reading_history (
    user_id UUID,
    blog_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, blog_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);

COMMIT;