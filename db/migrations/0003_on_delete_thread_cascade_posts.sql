ALTER TABLE posts
DROP CONSTRAINT posts_thread_id_fkey
;

ALTER TABLE posts
ADD CONSTRAINT posts_thread_id_fkey
FOREIGN KEY (thread_id)
REFERENCES threads(id)
ON DELETE CASCADE
;