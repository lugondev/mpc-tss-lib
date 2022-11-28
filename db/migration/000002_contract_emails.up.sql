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
