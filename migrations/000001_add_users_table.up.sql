CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    email varchar(200) unique not null,
    password varchar(200) not null
    
);