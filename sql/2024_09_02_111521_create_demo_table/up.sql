
create table demos (
    id number(19, 0) not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint demos_id_pk primary key (id)
);

create sequence demos_id_seq minvalue 1 start with 1 increment by 1;

create trigger demos_id_trg 
before insert on demos 
for each row 
begin 
if :new.ID is null then select demos_id_seq.nextval into :new.ID from dual; 
end if;
end;
/
