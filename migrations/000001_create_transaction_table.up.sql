BEGIN;

CREATE TABLE IF NOT EXISTS "transaction" (
    id serial not null unique,
    user_id integer not null,
    user_email varchar(254) not null,
    amount bigint not null,
    currency char(3) not null CHECK (length(currency) = 3),
    status varchar(8) not null default 'NEW' CHECK (status IN ('NEW', 'ERROR', 'SUCCESS', 'FAILED', 'CANCELED')),
    created_at timestamp with time zone default now()::timestamptz,
    updated_at timestamp with time zone default now()::timestamptz
);

COMMIT;