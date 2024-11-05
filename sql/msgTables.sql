create table chat_group_members
(
    group_id   bigint      not null,
    user_id    bigint      not null,
    created_at datetime(3) null,
    updated_at datetime(3) null,
    primary key (group_id, user_id)
);

create table chat_groups
(
    gid        bigint auto_increment        primary key,
    name       varchar(64)  null,
    avatar     varchar(256) null,
    owner_id   bigint       not null,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    deleted_at datetime(3)  null
);

create index idx_chat_groups_deleted_at
    on chat_groups (deleted_at);

create index idx_chat_groups_owner_id
    on chat_groups (owner_id);

create table group_messages
(
    msg_id    bigint auto_increment        primary key,
    sender_id bigint      not null,
    group_id  bigint      not null,
    content   text        not null,
    timestamp datetime(3) not null
);

create index idx_group_messages_group_id
    on group_messages (group_id);

create index idx_group_messages_sender_id
    on group_messages (sender_id);

create table messages
(
    msg_id      bigint auto_increment        primary key,
    sender_id   bigint      not null,
    receiver_id bigint      not null,
    content     text        not null,
    timestamp   datetime(3) not null
);

create index idx_messages_receiver_id
    on messages (receiver_id);

create index idx_messages_sender_id
    on messages (sender_id);



