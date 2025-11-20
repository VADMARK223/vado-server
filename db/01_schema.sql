CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(100) UNIQUE                                        NOT NULL,
    password   VARCHAR(255)                                               NOT NULL,
    email      VARCHAR(255) UNIQUE,
    role       varchar(20) CHECK (role IN ('user', 'moderator', 'admin')) NOT NULL,
    color      VARCHAR(7) CHECK (color ~ '^#[0-9A-Fa-f]{6}$')             NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
comment on table users is 'Таблица пользователей';
comment on column users.color is 'Цвет пользователя в HEX (#RRGGBB)';
comment on column users.role is 'Роль пользователя';

create table if not exists tasks
(
    id          serial primary key,
    name        varchar(255)                        not null,
    description text,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    completed   boolean   default false,
    updated_at  timestamp default CURRENT_TIMESTAMP not null,

    -- внешний ключ на таблицу users
    user_id     int                                 not null,
    constraint fk_tasks_users_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);
comment on table tasks is 'Таблица задач';