create or replace view v_session as 
    select 
        s.session,
        s.csrf_token,
        u.id as user_id,
        u.phone as user_phone,
        u.email as user_email,
        fn_trim((coalesce(u.sname, ''::text) || ' '::text) || coalesce(u.name, ''::text)) as user_name,
        s.life_time
    from t_session s
    join t_user u on 
        u.id = s.user_id
    where s.life_time > current_timestamp;
alter table v_session owner to postgres;
grant all on table v_session to postgres;
grant select on table v_session to api_user;
comment on view v_session is 'Сессии пользователей';
