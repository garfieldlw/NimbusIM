create table conversation
(
    id                bigint not null
        primary key,
    owner             bigint not null,
    conversation_id   bigint not null,
    conversation_type int    not null,
    user_id_1         bigint not null,
    user_id_2         bigint not null,
    status            int    not null,
    create_time       bigint not null,
    update_time       bigint not null,
    constraint conversation_owner_conversation_id_uindex
        unique (owner, conversation_id)
);

create table conversation_group
(
    id          bigint        not null
        primary key,
    user_id     bigint        not null,
    group_name  varchar(255)  not null,
    status      int default 1 not null,
    create_time bigint        not null,
    update_time bigint        not null
);

create table message
(
    id              bigint not null
        primary key,
    conversation_id bigint not null,
    quote           bigint not null,
    send_id         bigint not null,
    receive_id      bigint not null,
    from_type       int    not null,
    session_type    int    not null,
    content_type    int    not null,
    content_status  int    not null,
    content         json   not null,
    send_time       bigint not null,
    privacy         int    not null,
    status          int    not null,
    create_time     bigint not null,
    update_time     bigint not null
);

create table user
(
    id          bigint        not null
        primary key,
    email       varchar(255)  not null,
    status      int default 1 not null,
    create_time bigint        not null,
    update_time bigint        not null,
    constraint user_email_uindex
        unique (email)
);

