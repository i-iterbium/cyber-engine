create or replace function fn_trim(arg_str text)
returns text AS
$body$
begin
    return trim(chr(9) || chr(10) || chr(13) || chr(32) from arg_str);
end;
$body$
language plpgsql volatile cost 100;
alter function fn_trim(text) owner to postgres;
grant execute on function fn_trim(text) to postgres;
grant execute on function fn_trim(text) to api_user;
revoke all on function fn_trim(text) from public;
comment on function fn_trim(text) is 'Обрезает " ", \\t, \\r, \\n по краям строк';
