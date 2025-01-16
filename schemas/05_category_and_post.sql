CREATE TABLE IF NOT EXISTS category_and_post(
    category_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
    UNIQUE(category_id, post_id)
);