-- function: fn_check_user_session(text)

-- drop function fn_check_user_session(text);

create or replace function fn_check_user_session(arg_session text, arg_csrf_token text)
returns table (
    user_id int,
    user_role text
) as
$body$
declare
    v_session record;
    v_status text;
    v_role text;
begin
    perform fn_exception_if(arg_session is null, 'Не переданы данные сессии пользователя');
    perform fn_exception_if(arg_csrf_token is null, 'Не переданы данные сессии пользователя');

    select
        s.user_id,
        s.life_time
    from t_session s
    where 
        s.session = arg_session and
        s.csrf_token = arg_csrf_token
    into v_session;

    perform fn_exception_if(v_session.user_id is null, 'Пользователь не авторизован');

    select 
	    s.code
    from t_user u
	join t_user_status s on s.id = u.status_id 
    where
        u.id = v_session.user_id
	into v_status;
	
	perform fn_exception_if(v_status = 'new', 'Пользователь не подтвержден');
	perform fn_exception_if(v_status = 'blocked', 'Пользователь заблокирован');

    perform fn_exception_if(v_session.life_time < current_timestamp, 'Время жизни сессии пользователя истекло');

   -- select r.code
   -- from t_user_role r
   -- where r.id = v_session.role_id
   -- into v_role;

    --perform fn_exception_if(v_role is null, 'Не найдена роль пользователя');

    return query
    select
        v_session.user_id,
        v_role;
end;
$body$
    language plpgsql volatile;
alter function fn_check_user_session(text, text) owner to postgres;
grant execute on function fn_check_user_session(text, text) to postgres;
grant execute on function fn_check_user_session(text, text) to api_user;
comment on function fn_check_user_session(text) is 'Сессии пользователей. Проверка';
