INSERT into "league" (title)
VALUES
('Premier League');
INSERT into "league" (title)
VALUES
('Serie A');

INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at, league_id)
VALUES
('Hello World', 'hello-world', 'hello, world', 'Hello, World!', 1.00, 'over 1.5', now() + interval '1 day', 1);

INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at, league_id)
VALUES
('Hello World 2', 'hello-world-2', 'hello, world', 'Hello, World 2!', '2.00', 'under 4.5', now() + interval '2 day', 2);

INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at, is_featured, league_id)
VALUES
('Hello World 3', 'hello-world-3', 'hello, world', 'Hello, World 3!', '3.00','over 1.5', now() + interval '3 day', true, 2);