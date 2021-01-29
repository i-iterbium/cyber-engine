-- function: fn_session_by_phone_ins(bigint, text)

-- drop function fn_session_by_phone_ins(bigint, text);

create or replace function fn_session_by_phone_ins(
    arg_phone bigint,
    arg_password text
) returns table (
    user_id int,
    session varchar,
    csrf_token varchar,
    life_time timestamp without time zone
) as
$body$
declare
    С_DURATION interval = '30 days';

    v_password_salt text;
    v_password_hash text;
    v_user record;
    v_status text;
    v_session_id int;
begin
    perform fn_exception_if(arg_phone is null, 'Не передан номер телефона пользователя');
    perform fn_exception_if(arg_password is null, 'Не передан пароль пользователя');

    select
        u.id,
        u.password_hash,
        u.password_salt,
        u.status_id
    from t_user u
    where u.phone = arg_phone
    into v_user;

    perform fn_exception_if(v_user.id is null, 'Пользователь с номером телефона %s не зарегистрирован', arg_phone::text);
    perform fn_check_user(v_user.id);
    
    select 
	    s.code
	from t_user_status s 
	where s.id = v_user.status_id
	into v_status;
	
	perform fn_exception_if(v_status = 'new', 'Пользователь с номером телефона %s не подтвержден', arg_phone::text);
	perform fn_exception_if(v_status = 'blocked', 'Пользователь с номером телефона %s заблокирован', arg_phone::text);

    v_password_hash = pgcrypto.crypt(arg_password, v_user.password_salt);
    perform fn_exception_if(v_password_hash != v_user.password_hash, 'Передан неверный пароль');

    select s.id
    from t_session s
    where s.user_id = v_user.id
    into v_session_id;

    if not found then
        insert into t_session as s (
            user_id,
            session,
            csrf_token,
            life_time
        ) values (
            v_user.id,
            fn_generate_session_hash(),
            md5(random()::text),
            current_timestamp + С_DURATION
        ) returning s.id into v_session_id;
    else
        update t_session s set
            session     = fn_generate_session_hash(),
            csrf_token  = md5(random()::text),
            life_time   = current_timestamp + С_DURATION
        where s.user_id = v_user.id;
    end if;

    return query
    select
        s.user_id,
        s.session,
        s.csrf_token,
        s.life_time
    from t_session s
    where s.id = v_session_id;
end;
$body$
    language plpgsql volatile;
alter function fn_session_by_phone_ins(bigint, text) owner to postgres;
grant execute on function fn_session_by_phone_ins(bigint, text) to postgres;
grant execute on function fn_session_by_phone_ins(bigint, text) to api_user;
comment on function fn_session_by_phone_ins(bigint, text) is 'Сессии пользователей. Создание по номеру телефона';
