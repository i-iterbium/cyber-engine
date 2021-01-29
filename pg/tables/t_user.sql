create table t_user (
    id              serial  not null  primary key,
    password_hash   text    not null  unique,
    password_salt   text    not null,
    phone           bigint  not null  unique,
    email           text    not null  unique,
    birthday        timestamp without time zone,
    name            text,
    sname           text,
    pname           text,
    status_id       int     not null  references t_user_status (id)
);
comment on table t_user is 'Пользователи';
comment on column t_user.id is 'Идентификатор';
comment on column t_user.password_hash is 'Хэш пароля';
comment on column t_user.password_salt is 'Соль пароля';
comment on column t_user.phone is 'Номер телефона';
comment on column t_user.email is 'Email';
comment on column t_user.birthday is 'Дата рождения';
comment on column t_user.name is 'Имя';
comment on column t_user.sname is 'Фамилия';
comment on column t_user.pname is 'Отчество';
comment on column t_user.status_id is 'Идентификатор статуса';

alter table t_user owner to postgres;
grant all on table t_user to api_user;