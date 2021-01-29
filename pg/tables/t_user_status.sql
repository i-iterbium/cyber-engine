create table t_user_status (
    id      serial  not null  primary key,
    code    text    not null  unique,
    name    text    not null
);
comment on table t_user_status is 'Статусы пользователей';
comment on column t_user_status.id is 'Идентификатор';
comment on column t_user_status.code is 'Символьный код';
comment on column t_user_status.name is 'Наименование';

alter table t_user_status owner to postgres;
grant all on table t_user_status to api_user;

insert into t_user_status (
    code, name
) values
(
    'new', 'Новый'
),
(
    'active', 'Активный'
),
(
    'blocked', 'Заблокированный'
);
