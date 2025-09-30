-- +goose Up
ALTER TABLE comments ADD COLUMN parent_comment_id UUID REFERENCES comments(comment_id);
ALTER TABLE comments ADD COLUMN reply_to_user_id UUID; -- ID пользователя, которому отвечают

-- Индекс для быстрого поиска ответов на комментарий
CREATE INDEX idx_comments_parent ON comments(parent_comment_id) WHERE parent_comment_id IS NOT NULL;

-- Индекс для быстрого получения комментариев к посту
CREATE INDEX idx_comments_post_created ON comments(post_id, created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_comments_post_created;
DROP INDEX IF EXISTS idx_comments_parent;
ALTER TABLE comments DROP COLUMN IF EXISTS reply_to_user_id;
ALTER TABLE comments DROP COLUMN IF EXISTS parent_comment_id;