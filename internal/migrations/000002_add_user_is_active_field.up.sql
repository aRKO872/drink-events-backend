ALTER TABLE users DROP COLUMN IF EXISTS "is_active";
ALTER TABLE users ADD column "is_active" boolean default true;