create table if not exists request (
    id            serial not null primary key,
    method        varchar(16)   not null,
    scheme        varchar(16)   not null,
    host          varchar(128)  not null,
    path          varchar(1024) not null,
    query_string  text          default '',
    post_params   text          default '',
    cookies       text          default '',
    headers       text          default '',
    body          bytea         default ''
);

create table if not exists response (
    id           serial not null primary key,
    request_id   serial references request (id) not null,
    code         int                            not null,
    cookies      text                           default '',
    headers      text                           default '',
    body         bytea                          default ''
);
