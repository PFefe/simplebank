-- migrations/20240704_add_users_and_modify_accounts.down.sql

-- Step 1: Remove foreign key constraints for entries and transfers tables
ALTER TABLE "entries"
    DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";
ALTER TABLE "transfers"
    DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE "transfers"
    DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";

-- Step 2: Drop the unique constraint on the accounts table
ALTER TABLE "accounts"
    DROP CONSTRAINT IF EXISTS "accounts_owner_currency_key";

-- Step 3: Revert the owner column in the accounts table
ALTER TABLE "accounts"
    ADD COLUMN "new_owner" varchar;
UPDATE "accounts"
SET "new_owner" = "owner";
ALTER TABLE "accounts"
    DROP COLUMN "owner";
ALTER TABLE "accounts"
    RENAME COLUMN "new_owner" TO "owner";

-- Step 4: Drop the users table
DROP TABLE IF EXISTS "users";

-- Step 5: Re-add the original foreign key constraints for entries and transfers tables
ALTER TABLE "entries"
    ADD CONSTRAINT "entries_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers"
    ADD CONSTRAINT "transfers_from_account_id_fkey" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers"
    ADD CONSTRAINT "transfers_to_account_id_fkey" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
