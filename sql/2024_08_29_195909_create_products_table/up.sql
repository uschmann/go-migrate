whenever sqlerror exit sql.sqlcode;

create table products (
    id number(19, 0) not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint products_id_pk primary key (id)
);

create sequence products_id_seq minvalue 1 start with 1 increment by 1;

create trigger products_id_trg 
before insert on products 
for each row 
begin 
if :new.ID is null then select products_id_seq.nextval into :new.ID from dual; 
end if;
end;
/
