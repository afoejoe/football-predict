CREATE TABLE "blog" (
    "id" bigserial PRIMARY KEY,
    "title" varchar(255) NOT NULL,
    "slug" varchar(255) NOT NULL,
    "body" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
)

