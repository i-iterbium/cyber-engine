-- function: fn_session_upd(text)

-- drop function fn_session_upd(text);

create or replace function fn_session_upd(arg_session text, arg_csrf_token text)
returns table (
    user_id int,
    session text,
    csrf_token text,
    life_time text
) as
$body$
declare
    c_duration interval = '30 days';
    v_user_id int;
begin
    perform fn_exception_if(arg_session is null, 'Не передан токен сессии пользователя');
    perform fn_exception_if(arg_csrf_token is null, 'Не передан crsf-токен сессии пользователя');

    select
        s.user_id
    from v_session s
    where 
        s.session = arg_session and
        s.csrf_token = arg_csrf_token
    into v_user_id;

    perform fn_exception_if(v_user_id is null, 'Пользователь с сессией %s не существует', arg_session::text);

    update t_session s set
        life_time = current_timestamp + С_DURATION
    where s.user_id = v_user_id;

    return query
    select
        s.user_id,
        s.session,
        s.csrf_token,
        s.life_time
    from v_session s
    where s.user_id = v_user_id;
end;
$body$
    language plpgsql volatile;
alter function fn_session_upd(text, text) owner to postgres;
grant execute on function fn_session_upd(text, text) to postgres;
grant execute on function fn_session_upd(text, text) to api_user;
comment on function fn_session_upd(text, text) is 'Сессии пользователей. Продление срока действия';
