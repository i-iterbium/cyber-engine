create or replace function fn_check_password(arg_password text)
returns text AS
$body$
declare
    v_password_errors text[];
    v_alphas text[] = regexp_matches(arg_password, '[[:alpha:]]');
    v_digits text[] = regexp_matches(arg_password, '[[:digit:]]');
begin
    if length(coalesce(arg_password, '')) < 6 then
        v_password_errors = array_append(v_password_errors, 'минимум 6 символов');
    end if;

    v_alphas = regexp_matches(arg_password, '[[:alpha:]]');
    if v_alphas is null or array_length(v_alphas, 1) = 0 then
        v_password_errors = array_append(v_password_errors, 'латинские буквы');
    end if;

    v_digits = regexp_matches(arg_password, '[[:digit:]]');
    if v_digits is null or array_length(v_digits, 1) = 0 then
        v_password_errors = array_append(v_password_errors, 'цифры');
    end if;

    perform fn_exception_if(array_length(v_password_errors, 1) > 0, 'Некорректный пароль, пароль должен содержать %s', array_to_string(v_password_errors, ', '));

    return arg_password;
end;
$body$
language plpgsql volatile cost 100;
alter function fn_check_password(text) owner to postgres;
grant execute on function fn_check_password(text) to postgres;
grant execute on function fn_check_password(text) to api_user;
revoke all on function fn_check_password(text) from public;
comment on function fn_check_password(text) is 'Пароль пользователя. Проверка';
