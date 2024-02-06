CREATE TABLE "students" (
  "student_id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar,
  "phone_number" varchar,
  "address" text,
  "college_id" bigint,
  "funnel_id" bigint,
  "hourly_fee" float,
  "notes" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "colleges" (
  "college_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "funnels" (
  "funnel_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "lessons" (
  "lesson_id" bigserial PRIMARY KEY,
  "lesson_datetime" timestamptz NOT NULL,
  "duration" bigint NOT NULL,
  "location_id" bigint NOT NULL,
  "subject_id" bigint NOT NULL,
  "notes" text
);

CREATE TABLE "lesson_locations" (
  "location_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "lesson_subjects" (
  "subject_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "invoices" (
  "invoice_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "lesson_id" bigint NOT NULL,
  "hourly_fee" float NOT NULL,
  "duration" bigint NOT NULL,
  "discount" float NOT NULL,
  "amount" float NOT NULL,
  "notes" text,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "receipts" (
  "receipt_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "receipt_datetime" timestamptz NOT NULL,
  "amount" float NOT NULL,
  "notes" text
);

CREATE TABLE "payments" (
  "payment_id" bigserial PRIMARY KEY,
  "receipt_id" bigint NOT NULL,
  "payment_datetime" timestamptz NOT NULL,
  "amount" float NOT NULL,
  "payment_method_id" bigint NOT NULL
);

CREATE TABLE "payment_methods" (
  "payment_method_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE INDEX ON "students" ("first_name");

CREATE INDEX ON "students" ("last_name");

CREATE INDEX ON "students" ("first_name", "last_name");

CREATE INDEX ON "students" ("college_id");

CREATE INDEX ON "students" ("funnel_id");

CREATE INDEX ON "colleges" ("name");

CREATE INDEX ON "funnels" ("name");

CREATE INDEX ON "lessons" ("lesson_datetime");

CREATE INDEX ON "lessons" ("location_id");

CREATE INDEX ON "lessons" ("location_id", "lesson_datetime");

CREATE INDEX ON "lessons" ("subject_id");

CREATE INDEX ON "lessons" ("subject_id", "lesson_datetime");

CREATE INDEX ON "lesson_locations" ("name");

CREATE INDEX ON "lesson_subjects" ("name");

CREATE INDEX ON "invoices" ("student_id");

CREATE INDEX ON "invoices" ("lesson_id");

CREATE INDEX ON "invoices" ("student_id", "lesson_id");

CREATE INDEX ON "receipts" ("student_id");

CREATE INDEX ON "receipts" ("receipt_datetime");

CREATE INDEX ON "receipts" ("receipt_datetime", "student_id");

CREATE INDEX ON "payments" ("receipt_id");

CREATE INDEX ON "payments" ("payment_datetime");

CREATE INDEX ON "payment_methods" ("name");

COMMENT ON COLUMN "students"."hourly_fee" IS 'hourly fee for the student';

COMMENT ON COLUMN "lessons"."duration" IS 'lesson duration in minutes';

COMMENT ON COLUMN "invoices"."hourly_fee" IS 'hourly fee for the lesson';

COMMENT ON COLUMN "invoices"."duration" IS 'lesson duration in minutes';

COMMENT ON COLUMN "invoices"."amount" IS 'total amount based on lesson duration, hourly fee and discount';

COMMENT ON COLUMN "receipts"."amount" IS 'total amount of all payments';

ALTER TABLE "students" ADD FOREIGN KEY ("college_id") REFERENCES "colleges" ("college_id");

ALTER TABLE "students" ADD FOREIGN KEY ("funnel_id") REFERENCES "funnels" ("funnel_id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("location_id") REFERENCES "lesson_locations" ("location_id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("subject_id") REFERENCES "lesson_subjects" ("subject_id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("student_id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("lesson_id");

ALTER TABLE "receipts" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("student_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("receipt_id") REFERENCES "receipts" ("receipt_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("payment_method_id") REFERENCES "payment_methods" ("payment_method_id");
