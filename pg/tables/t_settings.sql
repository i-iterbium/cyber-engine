-- table: t_settings

-- drop table t_settings;

create table t_settings (
    id          serial      not null  primary key,
    name        text        not null  unique,
    description text        not null,
    value       jsonb       not null
);

comment on table t_settings is 'Параметры';
comment on column t_settings.id is 'Идентификатор';
comment on column t_settings.name is 'Название';
comment on column t_settings.description is 'Описание';
comment on column t_settings.value is 'Значение';

alter table t_settings owner to postgres;
grant all on table t_settings to api_user;
