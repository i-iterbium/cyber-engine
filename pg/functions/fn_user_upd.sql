create or replace function fn_user_upd(
    in arg_user_id integer,
    in arg_name text default null::text,
    in arg_sname text default null::text,
    in arg_pname text default null::text,
    in arg_birthday timestamp without time zone default null::timestamp without time zone,
    in arg_email text default null::text)
returns table(
    id integer, 
    phone bigint, 
    name text, 
    sname text, 
    pname text, 
    email text, 
    birthday timestamp without time zone) as
$body$
begin
    update t_user as u set
        name               = nullif(fn_trim(arg_name), ''),
        sname              = nullif(fn_trim(arg_sname), ''),
        pname              = nullif(fn_trim(arg_pname), ''),
        birthday           = arg_birthday,
        email              = coalesce(fn_trim(arg_email), u.email)
    where u.id = arg_user_id;

    perform fn_exception_if(not found, 'Пользователь с кодом %L не существует', arg_user_id::text);

    return query
    select
        u.id,
        u.phone,
        u.name,
        u.sname,
        u.pname,
        u.email,
        u.birthday
    from v_user u
    where u.id = arg_user_id;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_user_upd(integer, text, text, text, timestamp without time zone, text) owner to postgres;
grant execute on function fn_user_upd(integer, text, text, text, timestamp without time zone, text) to postgres;
grant execute on function fn_user_upd(integer, text, text, text, timestamp without time zone, text) to api_user;
revoke all on function fn_user_upd(integer, text, text, text, timestamp without time zone, text) from public;
comment on function fn_user_upd(integer, text, text, text, timestamp without time zone, text) is ' Пользователи. Изменение';
