
-- Link all posts to the General category
INSERT OR IGNORE INTO category_and_post (category_id, post_id)
SELECT 
    (SELECT id FROM category WHERE name = 'General'),
    id 
FROM post 
ORDER BY id ASC;