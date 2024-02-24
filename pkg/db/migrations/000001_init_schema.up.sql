CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT (false),
  "phone" varchar UNIQUE,
  "hashed_password" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
