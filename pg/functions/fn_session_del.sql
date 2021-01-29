-- function: fn_session_del(text)

-- drop function fn_session_del(text);

create or replace function fn_session_del(arg_session text, arg_csrf_token text)
returns void as
$body$
declare
    v_user_id int;
begin
    perform fn_exception_if(arg_session is null, 'Не передан токен сессии пользователя');
    perform fn_exception_if(arg_csrf_token is null, 'Не передан crsf-токен сессии пользователя');

    select
        s.user_id
    from t_session s
    where
        s.session = arg_session and
        s.csrf_token = arg_csrf_token
    into v_user_id;

    perform fn_exception_if(v_user_id is null, 'Пользователь с сессией %s не существует', arg_session::text);

    delete from t_session s
    where s.user_id = v_user_id;
end;
$body$
    language plpgsql volatile;
alter function fn_session_del(text) owner to postgres;
grant execute on function fn_session_del(text) to postgres;
grant execute on function fn_session_del(text) to api_user;
comment on function fn_session_del(text) is 'Сессии пользователей. Удаление';
