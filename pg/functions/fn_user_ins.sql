create or replace function fn_user_ins(
    in arg_phone bigint,
    in arg_password text,
    in arg_name text default null::text,
    in arg_sname text default null::text,
    in arg_pname text default null::text,
    in arg_email text default null::text,
    in arg_birthday timestamp without time zone default null::timestamp without time zone)
returns table(
    id integer, 
    phone bigint, 
    name text, 
    sname text, 
    pname text, 
    email text, 
    birthday timestamp without time zone,
    status text) AS
$body$
declare
    v_salt text;
    v_user_id int;
    v_status_new_id int;
begin
    perform fn_exception_if(
        exists(
            select 1
            from t_user u
            where u.phone = arg_phone
        ),
        'Пользователь с номером телефона %s уже зарегистрирован',
        arg_phone::text
    );

    v_salt = pgcrypto.gen_salt('bf');

    select
        us.id
    from t_user_status us
    where us.code = 'new'
    into v_status_new_id;

    insert into t_user as c (
        name,
        sname,
        pname,
        phone,
        password_hash,
        password_salt,
        email,
        birthday,
        status_id
    ) values (
        nullif(fn_trim(arg_name), ''),
        nullif(fn_trim(arg_sname), ''),
        nullif(fn_trim(arg_pname), ''),
        arg_phone,
        pgcrypto.crypt(fn_check_password(arg_password), v_salt),
        v_salt,
        nullif(fn_trim(arg_email), ''),
        arg_birthday,
        v_status_new_id
    ) returning c.id into v_user_id;

    return query
    select
        u.id,
        u.phone,
        u.name,
        u.sname,
        u.pname,
        u.email,
        u.birthday,
        u.status
    from v_user u
    where u.id = v_user_id;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_user_ins(bigint, text, text, text, text, text, timestamp without time zone) owner to postgres;
grant execute on function fn_user_ins(bigint, text, text, text, text, text, timestamp without time zone) to postgres;
grant execute on function fn_user_ins(bigint, text, text, text, text, text, timestamp without time zone) to api_user;
revoke all on function fn_user_ins(bigint, text, text, text, text, text, timestamp without time zone) from public;
comment on function fn_user_ins(bigint, text, text, text, text, text, timestamp without time zone) is 'Пользователи. Создание';
