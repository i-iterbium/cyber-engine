create or replace function fn_exception_if(
    in arg_condition boolean,
    in arg_message text,
    variadic arg_params text[] default array[]::text[])
returns void as
$body$
begin
    if arg_condition then
        raise exception 'USER_ERROR %', format(arg_message, variadic arg_params);
    end if;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_exception_if(boolean, text, text[]) owner to postgres;
grant execute on function fn_exception_if(boolean, text, text[]) to postgres;
grant execute on function fn_exception_if(boolean, text, text[]) to api_user;
revoke all on function fn_exception_if(boolean, text, text[]) from public;
comment on function fn_exception_if(boolean, text, text[]) is 'Исключение по условию';