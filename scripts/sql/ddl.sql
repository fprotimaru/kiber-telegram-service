create table "telegram_users" (
    id serial primary key,
    chat_id integer,
    phone varchar,
    unique (chat_id, phone)
);