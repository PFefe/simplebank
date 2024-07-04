CREATE TABLE "accounts"
(
    "id"         bigserial PRIMARY KEY,
    "owner"      varchar     NOT NULL,
    "balance"    bigint      NOT NULL,
    "currency"   varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries"
(
    "id"         bigserial PRIMARY KEY,
    "account_id" bigserial   NOT NULL,
    "amount"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers"
(
    "id"              bigserial PRIMARY KEY,
    "from_account_id" bigserial   NOT NULL,
    "to_account_id"   bigserial   NOT NULL,
    "amount"          bigint      NOT NULL,
    "created_at"      timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative';

COMMENT ON COLUMN "transfers"."amount" IS 'can not be a negative';

ALTER TABLE "entries"
    ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers"
    ADD CONSTRAINT transfers_from_account_id_fkey FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers"
    ADD CONSTRAINT transfers_to_account_id_fkey FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
