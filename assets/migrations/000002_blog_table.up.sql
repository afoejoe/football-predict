CREATE TABLE "prediction" (
    "id" bigserial PRIMARY KEY,
    "title" text NOT NULL,
    "slug" text UNIQUE NOT NULL,
    "keywords" text NOT NULL,
    "fulltext_search"  tsvector NOT NULL GENERATED ALWAYS AS (to_tsvector('english', title || ' ' || keywords || ' ' || body)) STORED,
    "body" text NOT NULL,
    "coefficient" decimal(5, 2) NOT NULL,
    "scheduled_at" timestamptz NOT NULL DEFAULT (now()),
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
)

