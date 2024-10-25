BEGIN;

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    body TEXT,
    claps INT,
    replies INT,
    author_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,    
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comments_blog_id ON comments(blog_id);

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);


CREATE TABLE IF NOT EXISTS clapped_blogs (
    user_id UUID NOT NULL,
    blog_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    
    FOREIGN KEY (blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
    CONSTRAINT unique_clapped_blogs_user_id_blog_id UNIQUE (user_id, blog_id)
);


CREATE TABLE IF NOT EXISTS clapped_comments  (
    user_id UUID NOT NULL,
    comment_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    CONSTRAINT unique_clapped_comments_user_id_comment_id UNIQUE (user_id, comment_id)
);

COMMIT;