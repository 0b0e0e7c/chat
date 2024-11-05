create table users
(
    uid        bigint auto_increment
        primary key,
    username   varchar(64)  null,
    password   varchar(128) null,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    deleted_at datetime(3)  null,
    constraint idx_users_username
        unique (username)
);

create table user_profiles
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(3)                  null,
    updated_at datetime(3)                  null,
    deleted_at datetime(3)                  null,
    uid        bigint                       not null,
    avatar     varchar(256)                 null,
    nickname   varchar(64)                  null,
    status     tinyint unsigned default '1' null,
    email      varchar(128)                 null,
    constraint idx_user_profiles_uid
        unique (uid),
    constraint fk_users_profile
        foreign key (uid) references users (uid)
);

create index idx_user_profiles_deleted_at
    on user_profiles (deleted_at);

create index idx_users_deleted_at
    on users (deleted_at);

create index idx_users_uid
    on users (uid);

