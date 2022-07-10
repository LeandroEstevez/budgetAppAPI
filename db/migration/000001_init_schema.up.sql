CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "total_expenses" decimal NOT NULL DEFAULT 0,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar NOT NULL,
  "name" varchar UNIQUE NOT NULL,
  "due_date" date NOT NULL,
  "amount" decimal NOT NULL DEFAULT 0
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "entries" ("owner");

COMMENT ON COLUMN "entries"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
