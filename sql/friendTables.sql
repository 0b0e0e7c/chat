create table friends
(
    id           bigint unsigned auto_increment
        primary key,
    created_at   datetime(3)      null,
    updated_at   datetime(3)      null,
    deleted_at   datetime(3)      null,
    user_id      bigint           not null,
    friend_id    bigint           not null,
    status       tinyint unsigned not null,
    initiator_id bigint           not null
);

create index idx_friends_deleted_at
    on friends (deleted_at);

create index idx_friends_friend_id
    on friends (friend_id);

create index idx_friends_user_id
    on friends (user_id);

create table friend_groups
(
    id         bigint auto_increment
        primary key,
    uid        bigint      not null,
    name       varchar(64) not null,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null
);

create index idx_friend_groups_deleted_at
    on friend_groups (deleted_at);

create index idx_friend_groups_uid
    on friend_groups (uid);

create table friend_group_members
(
    group_id   bigint      not null,
    friend_id  bigint      not null,
    joined_at  datetime(3) null,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    deleted_at datetime(3) null,
    primary key (group_id, friend_id),
    constraint fk_friend_groups_group_members
        foreign key (group_id) references friend_groups (id)
);

create index idx_friend_group_members_deleted_at
    on friend_group_members (deleted_at);

