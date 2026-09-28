INSERT INTO subscribe_types (label, created_at, updated_at)
VALUES ('free', now(), now())
ON CONFLICT (label) DO NOTHING;
