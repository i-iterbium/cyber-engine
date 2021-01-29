create or replace function fn_check_user(arg_user_id integer)
returns void as
$body$
declare
    v_user record;
begin
    perform fn_exception_if(arg_user_id is null, 'Не передан код пользователя');

    select
        us.code = 'blocked' as banned
    from t_user u
    join t_user_status us on 
        u.status_id = us.id
    where u.id = arg_user_id 
    into v_user;

    if not v_user.banned then
        return;
    else
        perform fn_exception_if(v_user.banned is null, 'Пользователь с кодом %s не зарегистрирован', arg_user_id::text);
        perform fn_exception_if(v_user.banned, 'Пользователь заблокирован');
    end if;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_check_user(integer) owner to postgres;
grant execute on function fn_check_user(integer) to postgres;
grant execute on function fn_check_user(integer) to api_user;
revoke all on function fn_check_user(integer) from public;
comment on function fn_check_user(integer) is 'Проверка пользователя';
