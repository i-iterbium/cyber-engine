-- function: fn_user_confirmation_code_get(int)

-- drop function fn_user_confirmation_code_get(int);

create or replace function fn_user_confirmation_code_get(arg_user_id int)
returns table (
    code text,
    timeout int
) as
$body$
declare
    v_status text;
begin
    perform fn_check_user(arg_user_id);

    select 
        us.code
    from t_user u
    join t_user_status us on us.id = u.status_id
    where u.id = arg_user_id
    into v_status;

    perform fn_exception_if(v_status != 'new', 'Подтверждение для пользователя %s недоступно', arg_user_id::text);

    return query
    select
        c.code,
        c.timeout
    from fn_sms_code_create(arg_user_id, 6) c;
end;
$body$
    language plpgsql volatile;
alter function fn_user_confirmation_code_get(int) owner to postgres;
grant execute on function fn_user_confirmation_code_get(int) to postgres;
grant execute on function fn_user_confirmation_code_get(int) to api_user;
comment on function fn_user_confirmation_code_get(int) is 'Код подтверждения регистрации пользователей. Получение';
