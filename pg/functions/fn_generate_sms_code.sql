-- function: fn_generate_sms_code(integer)

-- drop function fn_generate_sms_code(integer);

create or replace function fn_generate_sms_code(arg_number_of_digits int default 6)
  returns text as
$body$
declare
    NUMBER_OF_DIGITS int = coalesce(arg_number_of_digits, 6);
    v_code text = '';
begin
    for i in 1..NUMBER_OF_DIGITS loop
        v_code = v_code || round(random() * 9);
    end loop;
    
    return v_code;
    
--     if not exists(select 1 from t_sms_code c where c.code = v_code) then
--         return v_code;
--     end if;
--     return fn_generate_sms_code(arg_number_of_digits);
end;
$body$
    language plpgsql volatile;
alter function fn_generate_sms_code(int) owner to postgres;
grant execute on function fn_generate_sms_code(int) to postgres;
grant execute on function fn_generate_sms_code(iint) to api_user;
comment on function fn_generate_sms_code(int) is 'Генерация уникального кода для подтверждения sms';
