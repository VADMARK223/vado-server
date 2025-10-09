CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    email      VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
comment on table users is 'Таблица пользователей';

CREATE TABLE IF NOT EXISTS roles
(
    id   INT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL -- например: 'admin', 'user'
);
comment on table roles is 'Таблица ролей';

create table if not exists user_roles
(
    user_id integer references users (id) on delete cascade, -- Удаляем запись таблицы, если роль удалена
    role_id integer references roles (id) on delete cascade, -- Удаляем запись таблицы, если пользователь удален
    primary key (user_id, role_id)
);

comment on table user_roles is 'Таблица связи пользователей и ролей (многие-ко-многим)';

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