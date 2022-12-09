CREATE TABLE "chains"
(
    "id"          bigserial PRIMARY KEY,
    "name"        varchar     NOT NULL,
    "rpcs"      varchar[]    NOT NULL DEFAULT '{}',
    "chain_id"    bigint      NOT NULL UNIQUE,
    "updated_at"  timestamptz NOT NULL DEFAULT (now())
);
