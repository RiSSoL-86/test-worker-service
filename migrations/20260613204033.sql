-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" uuid NOT NULL,
  "title" text NOT NULL,
  "description" text NOT NULL,
  "customer_name" text NOT NULL,
  "customer_phone" text NOT NULL,
  "address" text NOT NULL,
  "priority" text NOT NULL,
  "status" text NOT NULL DEFAULT 'new',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
