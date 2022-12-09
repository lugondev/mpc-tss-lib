CREATE TYPE notification_status AS ENUM ('disable', 'enable');

CREATE TABLE "shares"
(
    "id"           bigserial PRIMARY KEY,
    "party_id"     varchar             NOT NULL UNIQUE,
    "pubkey"       varchar             NOT NULL UNIQUE,
    "data"         bytea               NOT NULL,
    "enable"       bool                NOT NULL DEFAULT true,
    "notification" notification_status NOT NULL DEFAULT 'enable',
    "address"      varchar(42)         NOT NULL UNIQUE CHECK (address ~* '^0x[a-fA-F0-9]{40}$'
) ,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "emails"
(
    "id"    bigserial PRIMARY KEY,
    "name"  varchar NOT NULL,
    "email" varchar NOT NULL UNIQUE CHECK (email ~* '^.+@.+\..+$'
) ,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "emails_notification"
(
    "id"         bigserial PRIMARY KEY,
    "email_id"   bigserial   NOT NULL,
    "pubkey"     varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT notification_email_unique UNIQUE (pubkey, email_id)
);

ALTER TABLE "emails_notification"
    ADD FOREIGN KEY ("email_id") REFERENCES "emails" ("id");
ALTER TABLE "emails_notification"
    ADD FOREIGN KEY ("pubkey") REFERENCES "shares" ("pubkey");

