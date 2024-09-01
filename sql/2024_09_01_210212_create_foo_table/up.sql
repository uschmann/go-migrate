whenever sqlerror exit sql.sqlcode;

create table foo (
    id number(19, 0) not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint foo_id_pk primary key (id)
);

create sequence foo_id_seq minvalue 1 start with 1 increment by 1;

create trigger foo_id_trg 
before insert on foo 
for each row 
begin 
if :new.ID is null then select foo_id_seq.nextval into :new.ID from dual; 
end if;
end;
/
