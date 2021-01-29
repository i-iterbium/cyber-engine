create or replace function fn_session_ins(
    in arg_phone bigint,
    in arg_password text)
returns table(
    session text, 
    csrf_token text, 
    user_id integer, 
    user_phone bigint, 
    life_time timestamp without time zone
) as
$body$
declare
    v_user record;
    v_password_hash text;
    c_duration interval = '30 days';
    v_session text;
begin
    select
        u.id,
        u.password_hash,
        u.password_salt
    from t_user u
    where u.phone = arg_phone
    into v_user;

    perform fn_exception_if(v_user is null, 'Пользователь с номером телефона %s не зарегистрирован', arg_phone::text);
    perform fn_check_user(v_user);

    v_password_hash = pgcrypto.crypt(arg_password, v_user.password_salt);
    perform fn_exception_if(v_password_hash != v_user.password_hash, 'Неверный пароль');

    v_session = fn_generate_session_hash();

    insert into t_session (
        session,
        csrf_token,
        user_id,
        expire_time
    ) values (
        v_session,
        md5(random()::text),
        arg_user_id,
        current_timestamp + c_duration
    );

    return query
    select
        s.session::text,
        s.csrf_token::text,
        s.user_id,
        s.user_phone,
        s.life_time
    from v_session s
    where s.session = v_session;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_session_ins(bigint, text) owner to postgres;
grant execute on function fn_session_ins(bigint, text) to postgres;
grant execute on function fn_session_ins(bigint, text) to api_user;
revoke all on function fn_session_ins(bigint, text) from public;
comment on function fn_session_ins(bigint, text) is 'Сессии пользователей. Создание';
