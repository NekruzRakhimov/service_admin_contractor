-- +goose Up
-- +goose StatementBegin
create table if not exists contractors_contractor
(
    id bigserial
    constraint contractors_contractor_pk
    primary key,
    resident boolean not null,
    bin varchar,
    name varchar,
    email varchar not null,
    block_date timestamp with time zone,
    status varchar default 'ACTIVE'::character varying not null,
    is_delete boolean default false not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contractors_contractor;
-- +goose StatementEnd
