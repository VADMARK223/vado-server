CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    login      VARCHAR(30) UNIQUE                                         NOT NULL,
    username   VARCHAR(30)                                                NOT NULL,
    password   VARCHAR(255)                                               NOT NULL,
    email      VARCHAR(100) UNIQUE                                        NOT NULL,
    role       varchar(20) CHECK (role IN ('user', 'moderator', 'admin')) NOT NULL,
    color      VARCHAR(7) CHECK (color ~ '^#[0-9A-Fa-f]{6}$')             NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
comment on table users is 'Таблица пользователей';
comment on column users.login is 'Логин пользователя';
comment on column users.username is 'Отображаемое имя пользователя';
comment on column users.email is 'Почта пользователя';
comment on column users.color is 'Цвет пользователя в HEX (#RRGGBB)';
comment on column users.role is 'Роль пользователя';

create table if not exists tasks
(
    id          serial primary key,
    name        varchar(255)                        NOT NULL,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP NOT NULL,
    completed   boolean   default false,
    updated_at  timestamp default CURRENT_TIMESTAMP NOT NULL,

    -- внешний ключ на таблицу users
    user_id     int                                 NOT NULL,
    constraint fk_tasks_users_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);
comment on table tasks is 'Таблица задач';

CREATE TABLE IF NOT EXISTS messages
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content      TEXT        NOT NULL,
    message_type VARCHAR(30) NOT NULL     DEFAULT 'text',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_user_id ON messages (user_id);
CREATE INDEX idx_messages_messages_type ON messages (message_type);