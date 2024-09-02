
create table {{ .name }} (
    id number(19, 0) not null,
    created_at timestamp null,
    updated_at timestamp null,
    constraint {{ .name }}_id_pk primary key (id)
);

create sequence {{ .name }}_id_seq minvalue 1 start with 1 increment by 1;

create trigger {{ .name }}_id_trg 
before insert on {{ .name }} 
for each row 
begin 
if :new.ID is null then select {{ .name }}_id_seq.nextval into :new.ID from dual; 
end if;
end;
/
