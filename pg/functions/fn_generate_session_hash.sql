create or replace function fn_generate_session_hash()
returns text as
$body$
declare
    v_hash text = md5(random()::text);
begin
    if not exists(select 1 from t_session s where s.session = v_hash) then
        return v_hash;
    end if;
    return fn_generate_session_hash();
end;
$body$
language plpgsql volatile cost 100;
alter function fn_generate_session_hash() owner to postgres;
grant execute on function fn_generate_session_hash() to postgres;
grant execute on function fn_generate_session_hash() to api_user;
revoke all on function fn_generate_session_hash() FROM public;
comment on function fn_generate_session_hash() IS 'Генерация уникального кода для сессии';
