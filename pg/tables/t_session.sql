create table t_session (
  id            serial                          not null  primary key,
  session       varchar(32)                     not null, 
  csrf_token    varchar(32)                     not null, 
  user_id       integer                         not null references t_user (id) match simple on update cascade on delete cascade, 
  life_time     timestamp without time zone     not null 
);
comment on table t_session is'Сессии пользователей';
comment on column t_session.id IS 'Идентификатор';
comment on column t_session.session IS 'Хэш сессиии';
comment on column t_session.csrf_token IS 'CSRF-токен';
comment on column t_session.user_id IS 'Идентификатор пользователя';
comment on column t_session.life_time IS 'Время истечения сессии';

alter table t_session owner to postgres;
grant all on table t_session to api_user;