
create table posts (
    id number(19, 0) not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint posts_id_pk primary key (id)
);

create sequence posts_id_seq minvalue 1 start with 1 increment by 1;

create trigger posts_id_trg 
before insert on posts 
for each row 
begin 
if :new.ID is null then select posts_id_seq.nextval into :new.ID from dual; 
end if;
end;
/
