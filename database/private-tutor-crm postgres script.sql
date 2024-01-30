CREATE TABLE "students" (
  "student_id" bigserial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar,
  "phone_number" varchar,
  "email_address" varchar,
  "college_id" integer,
  "funnel_id" integer,
  "notes" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "colleges" (
  "college_id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "funnels" (
  "funnel_id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "lessons" (
  "lesson_id" bigserial PRIMARY KEY,
  "held_on" timestamptz NOT NULL,
  "location_id" integer NOT NULL,
  "duration" time NOT NULL,
  "subject_id" integer NOT NULL,
  "notes" varchar
);

CREATE TABLE "lesson_locations" (
  "location_id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "lesson_subjects" (
  "subject_id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE TABLE "invoices" (
  "invoice_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "lesson_id" bigint NOT NULL,
  "lesson_duration" time NOT NULL,
  "hourly_amount" decimal(7, 2) NOT NULL,
  "total_amount" decimal(10, 2) NOT NULL,
  "notes" varchar
);

CREATE TABLE "receipts" (
  "receipt_id" bigserial PRIMARY KEY,
  "student_id" bigint NOT NULL,
  "received_at" timestamptz NOT NULL,
  "total_amount" decimal(10, 2) NOT NULL,
  "notes" varchar
);

CREATE TABLE "payments" (
  "payment_id" bigserial PRIMARY KEY,
  "receipt_id" bigint NOT NULL,
  "payed_at" timestamptz NOT NULL,
  "amount" DECIMAL(7, 2) NOT NULL,
  "payment_methods_id" integer NOT NULL
);

CREATE TABLE "payment_methods" (
  "payment_method_id" serial PRIMARY KEY,
  "name" varchar NOT NULL
);

CREATE INDEX ON "students" ("first_name");

CREATE INDEX ON "students" ("last_name");

CREATE INDEX ON "students" ("first_name", "last_name");

CREATE INDEX ON "students" ("college_id");

CREATE INDEX ON "students" ("funnel_id");

CREATE INDEX ON "lessons" ("location_id");

CREATE INDEX ON "lessons" ("subject_id");

CREATE INDEX ON "invoices" ("student_id");

CREATE INDEX ON "invoices" ("lesson_id");

CREATE INDEX ON "receipts" ("student_id");

CREATE INDEX ON "receipts" ("received_at");

CREATE INDEX ON "payments" ("receipt_id");

CREATE INDEX ON "payments" ("payed_at");

ALTER TABLE "students" ADD FOREIGN KEY ("college_id") REFERENCES "colleges" ("college_id");

ALTER TABLE "students" ADD FOREIGN KEY ("funnel_id") REFERENCES "funnels" ("funnel_id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("location_id") REFERENCES "lesson_locations" ("location_id");

ALTER TABLE "lessons" ADD FOREIGN KEY ("subject_id") REFERENCES "lesson_subjects" ("subject_id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("student_id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("lesson_id") REFERENCES "lessons" ("lesson_id");

ALTER TABLE "receipts" ADD FOREIGN KEY ("student_id") REFERENCES "students" ("student_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("receipt_id") REFERENCES "receipts" ("receipt_id");

ALTER TABLE "payments" ADD FOREIGN KEY ("payment_methods_id") REFERENCES "payment_methods" ("payment_method_id");
