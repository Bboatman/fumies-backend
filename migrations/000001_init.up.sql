BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "post_status" AS ENUM (
  'draft',
  'published',
  'private'
);

CREATE TABLE "metric" (
  "id"    uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "type"  varchar(80),
  "label" varchar(80),
  "t"     double precision,
  "f"     double precision,
  "g"     double precision,
  "m"     double precision
);

CREATE TABLE "user" (
  "id"        varchar(126) PRIMARY KEY,
  "username"  varchar(126),
  "email"     varchar(126)
);

CREATE TABLE "perfume" (
  "id"          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "name"        varchar(126),
  "house"       varchar(126),
  "url"         varchar(126),
  "is_empty"    boolean,
  "description" text,
  "user_id"     varchar(126),
  "created_at"  timestamp DEFAULT (now())
);

CREATE TABLE "perfume_metric" (
  "id"          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "note_id"     uuid,
  "perfume_id"  uuid,
  "created_at"  timestamp DEFAULT (now())
);

CREATE TABLE "review" (
  "id" uuid     PRIMARY KEY DEFAULT uuid_generate_v4(),
  "title"       varchar,
  "body"        text,
  "user_id"     varchar(126),
  "status"      post_status,
  "created_at"  timestamp DEFAULT (now())
);

CREATE TABLE "wear" (
  "id"          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "created_at"  timestamp DEFAULT (now()),
  "perfume_id"  uuid,
  "user_id"     varchar(126)
);

ALTER TABLE "perfume" 
  ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "perfume_metric" 
  ADD FOREIGN KEY ("note_id") REFERENCES "metric" ("id");

ALTER TABLE "perfume_metric" 
  ADD FOREIGN KEY ("perfume_id") REFERENCES "perfume" ("id");

ALTER TABLE "review" 
  ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "wear" 
  ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "wear" 
  ADD FOREIGN KEY ("perfume_id") REFERENCES "perfume" ("id");
COMMIT;