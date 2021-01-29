-- function: fn_sms_code_create(int, int)

-- drop function fn_sms_code_create(int, int);

create or replace function fn_sms_code_create(
    arg_user_id int,
    arg_code_length int
) returns table (
    code text,
    timeout int
) as
$body$
declare
    v_sms_code_lifetime int = (select p.value from t_settings p where p.name = 'smsCodeLifetime');
    v_code_limit_exceeded boolean;
    v_timeout int;
    v_code text;
    v_exist boolean;
    v_max_count int = (select p.value from t_settings p where p.name = 'limitSMSPerDay');
begin
    -- счетчик обнуляется, если последний код был получен в предыдущий день
    update t_sms_code c set
        count = 0
    where
        c.user_id = arg_user_id and
        creation_time::date < current_date;

    select
        true,
        c.count > v_max_count,
        case 
            when extract(epoch from c.creation_time + make_interval(secs => v_sms_code_lifetime) - current_timestamp)::int > 0 then extract(epoch from c.creation_time + make_interval(secs => v_sms_code_lifetime) - current_timestamp)::int
        end
    from t_sms_code c
    where
        c.user_id = arg_user_id
    into
        v_exist,
        v_code_limit_exceeded,
        v_timeout;

    if v_code_limit_exceeded then
        perform fn_user_block(v_user_id);

        return query
        select
            null::text,
            null::int;

        return;
    end if;

    if v_timeout is not null then
        return query
        select
            null::text,
            v_timeout;

        return;
    end if;

    select *
    from fn_generate_sms_code(arg_code_length)
    into v_code;

    if v_exist then
        update t_sms_code c set
            code            = v_code,
            creation_time   = current_timestamp,
            count = count + 1
        where c.user_id = arg_user_id;
    else
        insert into t_sms_code (
            user_id,
            code,
            count
        ) values (
            arg_user_id,
            v_code,
            1
        );
    end if;

    return query
    select
        v_code,
        extract(epoch FROM make_interval(secs => v_sms_code_lifetime))::int;
end;
$body$
    language plpgsql volatile;
alter function fn_sms_code_create(int, int) owner to postgres;
grant execute on function fn_sms_code_create(int, int) to postgres;
grant execute on function fn_sms_code_create(int, int) to api_user;
comment on function fn_sms_code_create(int, int) is 'SMS-коды. Создание';
