CREATE TABLE IF NOT EXISTS comment_vote(
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    vote INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(comment_id) REFERENCES comment(id) ON DELETE CASCADE,
    UNIQUE(user_id, comment_id)
);