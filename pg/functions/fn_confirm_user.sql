-- function: fn_confirm_user(int, text)

-- drop function fn_confirm_user(int, text);

create or replace function fn_confirm_user(
    arg_user_id int,
    arg_code text
) returns void as
$body$
declare
    v_code_exists boolean;
    v_code_active boolean;
    v_sms_code_lifetime int = (select p.value from t_settings p where p.name = 'smsCodeLifetime');
    v_status text;
begin
    perform fn_exception_if(arg_code is null, 'Не передан код подтверждения');
    perform fn_exception_if(arg_user_id is null, 'Не передан идентификатор пользователя');

    perform fn_check_user(arg_user_id);

    select 
        us.code
    from t_user u
    join t_user_status us on us.id = u.status_id
    where u.id = arg_user_id
    into v_status;

    perform fn_exception_if(v_status != 'new', 'Подтверждение для пользователя %s недоступно', arg_user_id::text);

    select 
        true,
        c.creation_time + make_interval(secs => v_sms_code_lifetime) > current_timestamp
    from t_sms_code c
    where 
        c.code = arg_code and
        c.user_id = arg_user_id
    into 
        v_code_exists, 
        v_code_active;

    perform fn_exception_if(v_code_exists is null, 'Передан некорректный код подтверждения');
    perform fn_exception_if(v_code_active != true, 'Время действия кода истекло. Пожалуйста, запросите новый код подтверждения');

    update t_user as u set
        status_id       = (select s.id from t_user_status s where s.code = 'active')
    where u.id = arg_user_id;

    delete from t_sms_code c
    where 
        c.code = arg_code and 
        c.user_id = arg_user_id;
end;
$body$
    language plpgsql volatile;
alter function fn_confirm_user(int, text) owner to postgres;
grant execute on function fn_confirm_user(int, text) to postgres;
grant execute on function fn_confirm_user(int, text) to api_user;
comment on function fn_confirm_user(int, text) is 'Подтверждение регистрации пользователя';
