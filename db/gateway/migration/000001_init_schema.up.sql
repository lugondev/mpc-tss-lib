CREATE TYPE requisition_status AS ENUM ('pending', 'failure', 'success');
CREATE TYPE requisition_type AS ENUM ('keygen', 'sign', 'reshare');

CREATE TABLE "requisitions"
(
    "id"          bigserial PRIMARY KEY,
    "requisition" varchar            NOT NULL,
    "pubkey"      varchar            NOT NULL DEFAULT '',
    "data"        bytea              NOT NULL,
    "reasons"     varchar            NOT NULL DEFAULT '',
    "username"    varchar            NOT NULL DEFAULT '',
    "tenant"      varchar            NOT NULL DEFAULT '',
    "retryTimes"  integer            NOT NULL DEFAULT 0,
    "type"        requisition_type   NOT NULL,
    "status"      requisition_status NOT NULL DEFAULT 'pending',
    "created_at"  timestamptz        NOT NULL DEFAULT (now()),
    "updated_at"  timestamptz        NOT NULL DEFAULT (now())
);
