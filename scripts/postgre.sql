CREATE USER postgre WITH password 'fyvaoldzh';
CREATE DATABASE qdago OWNER postgre;
GRANT ALL privileges ON DATABASE qdago TO postgre;

CREATE TABLE users
(
    id       bigserial primary key,
    login    varchar(60) not null,
    password varchar(60) not null,
    name     varchar(60) not null,
    email    varchar(60),
    birthday date,
    city     varchar(60),
    about    text,
    avatar   varchar(60) DEFAULT 'public/users/default.png'
);

CREATE TABLE events
(
    id          bigserial primary key,
    title       varchar(100) not null,
    place       varchar(100) not null,
    subway      varchar(60),
    street      varchar(60),
    description text         not null,
    category    varchar(60)  not null,
    start_date  timestamp    not null,
    end_date    timestamp    not null,
    image       varchar(60) default 'public/events/default.png'
);


CREATE TABLE tags
(
    id   serial primary key,
    name varchar(60) not null,
    CONSTRAINT unique_tag_names UNIQUE (name)
);


CREATE TABLE event_tag
(
    event_id bigint references events (id) on delete cascade,
    tag_id   int references tags (id) on delete cascade,
    CONSTRAINT unique_rows_event_tag UNIQUE (event_id, tag_id)
);

CREATE TABLE user_event
(
    user_id     bigint references users (id) on delete cascade,
    event_id    bigint references events (id) on delete cascade,
    is_planning boolean not null,
    CONSTRAINT unique_rows_user_event UNIQUE (user_id, event_id)
);

CREATE TABLE subscriptions
(
    subscriber_id    bigint references users (id) on delete cascade,
    subscribed_to_id bigint references users (id) on delete cascade,
    CONSTRAINT unique_subscribes UNIQUE (subscriber_id, subscribed_to_id)
);


CREATE TABLE user_preference
(
    user_id bigint references users (id) on delete cascade,
    show    int DEFAULT 0,
    movie   int DEFAULT 0,
    concert int DEFAULT 0
);


CREATE TABLE categories
(
    id   serial primary key,
    name varchar(60) not null,
    CONSTRAINT unique_categories_names UNIQUE (name)
);

CREATE TABLE actions_user_event
(
    user_id  bigint references users (id) on delete cascade,
    event_id bigint references events (id) on delete cascade,
    time     timestamp not null
);

CREATE TABLE actions_subscription
(
    subscriber_id bigint references users (id) on delete cascade,
    subscribed_to_id bigint references users (id) on delete cascade,
    time    timestamp not null
);

create table dialogues
(
    id serial primary key,
    user_1 bigint references users(id) on delete cascade,
    user_2 bigint references users(id) on delete cascade
);

create table messages
(
    id serial primary key,
    id_dialogue bigint references dialogues(id) on delete cascade,
    mes_from bigint references users(id) on delete cascade,
    mes_to bigint references users(id) on delete cascade,
    text text,
    date timestamp not null,
    redact bool default false,
    read bool default false
);
