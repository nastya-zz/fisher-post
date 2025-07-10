-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-------------------
-- Основные таблицы
-------------------
CREATE TABLE posts (
                       post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL,
                       description TEXT,
                       latitude DECIMAL(10, 7),
                       longitude DECIMAL(11, 7),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE media (
                       media_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
                       media_type VARCHAR(10) NOT NULL CHECK (media_type IN ('PHOTO', 'VIDEO')),
                       url VARCHAR(255) NOT NULL,
                       thumbnail_url VARCHAR(255),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-------------------
-- Взаимодействия
-------------------

CREATE TABLE likes (
                       like_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
                       user_id UUID NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                       UNIQUE (post_id, user_id) -- 1 лайк на пост от пользователя
);

CREATE TABLE comments (
                          comment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
                          user_id UUID NOT NULL,
                          content TEXT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-------------------
-- Метаданные рыбалки
-------------------

CREATE TABLE fish_types (
                            fish_id SERIAL PRIMARY KEY,
                            name VARCHAR(50) NOT NULL UNIQUE,
                            description TEXT
);

CREATE TABLE tackle_types (
                              tackle_id SERIAL PRIMARY KEY,
                              name VARCHAR(50) NOT NULL UNIQUE,
                              category VARCHAR(50) -- (спиннинг, поплавок, фидер и т.д.)
);

CREATE TABLE post_fish (
                           post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
                           fish_id INTEGER NOT NULL REFERENCES fish_types(fish_id) ON DELETE CASCADE,
                           PRIMARY KEY (post_id, fish_id)
);

CREATE TABLE post_tackle (
                             post_id UUID NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
                             tackle_id INTEGER NOT NULL REFERENCES tackle_types(tackle_id) ON DELETE CASCADE,
                             PRIMARY KEY (post_id, tackle_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
