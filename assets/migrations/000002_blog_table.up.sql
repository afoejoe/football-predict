CREATE TABLE "league" (
    "id" bigserial PRIMARY KEY,
    "title" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "prediction" (
    "id" bigserial PRIMARY KEY,
    "title" text NOT NULL,
    "slug" text UNIQUE NOT NULL,
    "keywords" text NOT NULL,
    "fulltext_search"  tsvector NOT NULL GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || keywords || ' ' || body)) STORED,
    "body" text NOT NULL,
    "odds" decimal(5, 2) NOT NULL,
    "prediction_type" text NOT NULL,
    "scheduled_at" timestamptz NOT NULL DEFAULT (now()),
    "is_featured" boolean NOT NULL DEFAULT false,
    "is_archived" boolean NOT NULL DEFAULT false,
    "campaigned" boolean NOT NULL DEFAULT false,
    "league_id" bigserial NOT NULL references "league" ("id"),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
)