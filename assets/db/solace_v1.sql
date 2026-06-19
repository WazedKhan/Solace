CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "name" varchar,
  "email" varchar UNIQUE,
  "password" varchar,
  "created_at" datetime,
  "updated_at" datetime
);

CREATE TABLE "habits" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "title" varchar,
  "image_url" varchar,
  "created_at" datetime,
  "updated_at" datetime,
  "current_streak" int DEFAULT 0,
  "last_checked_at" datetime
);

CREATE TABLE "habit_checking" (
  "id" uuid PRIMARY KEY,
  "habit_id" uuid,
  "checked_date" date,
  "created_at" datetime
);

CREATE TABLE "moods" (
  "id" uuid PRIMARY KEY,
  "label" varchar
);

CREATE TABLE "journal" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "mood_id" uuid,
  "status" varchar,
  "description" text,
  "image_url" varchar,
  "created_at" datetime,
  "updated_at" datetime
);

ALTER TABLE "habits" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "habit_checking" ADD FOREIGN KEY ("habit_id") REFERENCES "habits" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "journal" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "journal" ADD FOREIGN KEY ("mood_id") REFERENCES "moods" ("id") DEFERRABLE INITIALLY IMMEDIATE;
