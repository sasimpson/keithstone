drop table if exists users;
create table if not exists users (
    id integer primary key autoincrement,
    name text,
    password text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
insert into users (name, password) values ("admin", "password");

drop table if exists apikeys;
create table if not exists apikeys (
    id integer primary key autoincrement,
    key text UNIQUE,
    user_id integer,
    created_by integer,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

drop table if exists tokens;
create table if not exists tokens (
    id integer primary key autoincrement,
    token_id text UNIQUE,
    user_id integer,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    expiration timestamp
);


