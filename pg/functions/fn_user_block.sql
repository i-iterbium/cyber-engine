-- function: fn_user_block(int)

-- drop function fn_user_block(int);

create or replace function fn_user_block(arg_user_id int)
returns void as
$body$
begin
    perform fn_check_user(fn_user_block);

    update t_user as u set
        status_id = (select s.id from t_user_status s where s.code = 'blocked')
    where u.id = arg_user_id;

    delete from t_session s
    where s.user_id = arg_user_id;
end;
$body$
    language plpgsql volatile
    cost 100;
alter function fn_user_block(int) owner to postgres;
grant execute on function fn_user_block(int) to postgres;
grant execute on function fn_user_block(int) to api_user;
comment on function fn_user_block(int) is 'Пользователи. Блокировка';
