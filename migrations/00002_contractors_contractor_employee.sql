-- +goose Up
-- +goose StatementBegin
create table if not exists contractors_contractor_employee
(
    id bigserial
    constraint contractors_contractor_employee_pk
    primary key,
    contractor_id bigint
    constraint contractors_contractor_employee_contractors_contractor_id_fk
    references contractors_contractor,
    email varchar not null,
    full_name varchar,
    position varchar,
    block_date timestamp with time zone,
    status varchar default 'ACTIVE'::character varying,
    is_delete boolean default false not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contractors_contractor_employee;
-- +goose StatementEnd
