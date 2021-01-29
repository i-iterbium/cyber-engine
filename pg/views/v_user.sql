create or replace view v_user as 
    select 
        u.id,
        u.phone,
        nullif(fn_trim(u.name), ''::text) as name,
        nullif(fn_trim(u.sname), ''::text) as sname,
        nullif(fn_trim(u.pname), ''::text) as pname,
        u.email,
        u.birthday,
        us.code as status
    from t_user u
    join t_user_status us on us.id = u.status_id;
alter table v_user owner to postgres;
grant all on table v_user to postgres;
grant select on table v_user to api_user;
comment on view v_user is 'Пользователи';
