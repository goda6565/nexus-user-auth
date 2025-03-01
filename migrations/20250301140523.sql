-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "obj_id" uuid NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(255) NOT NULL,
  "username" character varying(255) NOT NULL,
  "avatar_url" character varying(255) NULL,
  "email_verified_at" timestamptz NULL,
  "last_login_at" timestamptz NULL,
  "role" character varying(50) NOT NULL DEFAULT 'user',
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create index "idx_users_obj_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_obj_id" ON "public"."users" ("obj_id");
