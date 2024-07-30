DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    parent_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
    pub_time INTEGER NOT NULL,
    content TEXT NOT NULL
);
