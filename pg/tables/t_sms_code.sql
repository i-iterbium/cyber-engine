-- table: t_sms_code

-- drop table t_sms_code;

create table t_sms_code (
    id            serial     not null  primary key,
    user_id       int        not null  references t_user (id),
    code          text       not null  unique,
    creation_time timestamp  not null  default current_timestamp,
    count         int        not null
);

comment on table t_sms_code is 'SMS-коды подтверждения';
comment on column t_sms_code.id is 'Идентификатор';
comment on column t_sms_code.user_id is 'Идентификатор пользователя';
comment on column t_sms_code.code is 'Символьный код';
comment on column t_sms_code.creation_time is 'Время создания';
comment on column t_sms_code.count is 'Счетчик запросов';

alter table t_sms_code owner to postgres;
grant all on table t_sms_code to api_user;
