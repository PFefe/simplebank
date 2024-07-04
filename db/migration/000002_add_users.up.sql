-- migrations/20240704_add_users_and_modify_accounts.up.sql

-- Step 1: Create the new users table
CREATE TABLE "users"
(
    "username"            varchar PRIMARY KEY,
    "hashed_password"     varchar        NOT NULL,
    "full_name"           varchar        NOT NULL,
    "email"               varchar UNIQUE NOT NULL,
    "password_changed_at" timestamptz    NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at"          timestamptz    NOT NULL DEFAULT now()
);

-- Step 2: Alter the accounts table to add a foreign key to users
ALTER TABLE "accounts"
    ADD COLUMN "new_owner" varchar;

-- Step 3: Migrate data from the old owner column to the new_owner column
UPDATE "accounts"
SET "new_owner" = "owner";

-- Step 4: Drop the old owner column
ALTER TABLE "accounts"
    DROP COLUMN "owner";

-- Step 5: Rename new_owner column to owner
ALTER TABLE "accounts"
    RENAME COLUMN "new_owner" TO "owner";

-- Step 6: Add the foreign key constraint
ALTER TABLE "accounts"
    ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

-- Step 7: Add unique constraint to the accounts table
ALTER TABLE "accounts"
    ADD CONSTRAINT "accounts_owner_currency_key" UNIQUE ("owner", "currency");

-- Step 8: Modify foreign key constraints for entries and transfers tables
ALTER TABLE "entries"
    DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";
ALTER TABLE "entries"
    ADD CONSTRAINT "entries_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers"
    DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE "transfers"
    ADD CONSTRAINT "transfers_from_account_id_fkey" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers"
    DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";
ALTER TABLE "transfers"
    ADD CONSTRAINT "transfers_to_account_id_fkey" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
