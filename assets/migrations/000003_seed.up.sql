INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at)
VALUES
('Hello World', 'hello-world', 'hello, world', 'Hello, World!', 1.00, 'over 1.5', now() + interval '1 day');

INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at)
VALUES
('Hello World 2', 'hello-world-2', 'hello, world', 'Hello, World 2!', '2.00', 'under 4.5', now() + interval '2 day');

INSERT into "prediction" (title, slug, keywords, body, odds, prediction_type, scheduled_at, is_featured)
VALUES
('Hello World 3', 'hello-world-3', 'hello, world', 'Hello, World 3!', '3.00','over 1.5', now() + interval '3 day', true);