CREATE TYPE notification_status AS ENUM ('disable', 'enable');

CREATE TABLE "shares"
(
    "id"           bigserial PRIMARY KEY,
    "pubkey"       varchar             NOT NULL UNIQUE,
    "data"         jsonb,
    "enable"       bool                NOT NULL DEFAULT true,
    "notification" notification_status NOT NULL DEFAULT 'enable',
    "address"      varchar(42)         NOT NULL UNIQUE CHECK (address ~* '^0x[a-fA-F0-9]{40}$'
) ,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
