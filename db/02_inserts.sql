insert into roles (id, name) values (1, 'admin');
insert into roles (id, name) values (2, 'moderator');
insert into roles (id, name) values (3, 'user');

insert into users (username, password, email, created_at) values ('1', '$2a$10$uKQGFGjps52djEN1yYUkvO5cUuELFbqZgxFOyxI6D6kjwh5Ne2W5m', 'user1@mail.ru', now());
insert into users (username, password, email, created_at) values ('2', '$2a$10$AHiUObaB7UdeslZP2WAd.uCbXu01LspUz7KiLMPOfze67NIYJEcPy', 'user2@mail.ru', now());