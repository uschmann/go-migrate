whenever sqlerror exit sql.sqlcode;

DROP TABLE {{ .name }};

DROP SEQUENCE {{ .name }}_id_seq;