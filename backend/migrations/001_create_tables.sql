create table users (
    id serial primary key,
    username varchar(50) unique not null,
    profile_picture text null,
    password_hash varchar(60) not null,
    created_at timestamp default current_timestamp
);

create table notes (
    id serial primary key,
    relationship_id int references relationships(id) on delete cascade,
    author_id int references users(id) on delete cascade,
    title varchar(255) not null,
    content text not null,
    position_x decimal(10, 2) default 0,
    position_y decimal(10, 2) default 0,
    color varchar(7) not null default '#FFFFFF',
    created_at timestamp default current_timestamp
);

create table relationships (
    id serial primary key,
    name text not null,
    picture text null,
    created_at timestamp default current_timestamp
);

create table relationship_members (
    relationship_id int references relationships(id) on delete cascade,
    user_id int references users(id) on delete cascade,
    primary key (relationship_id, user_id)
);

create table invites (
    id serial primary key,
    relationship_id int not null references relationships(id) on delete cascade,
    inviter_id int not null references users(id) on delete cascade,
    invitee_id int not null references users(id) on delete cascade,
    body text not null default 'be mine <3'
);