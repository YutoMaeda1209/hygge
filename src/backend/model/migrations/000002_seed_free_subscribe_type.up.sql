-- The id is fixed because the application refers to it as a constant (model.freeSubscribeTypeId).
INSERT INTO
    subscribe_types (id, label, created_at, updated_at)
VALUES
    (1, 'free', now(), now()) ON CONFLICT (id) DO NOTHING;

-- Advance the sequence past the explicit id so later inserts do not collide.
SELECT
    setval(
        'subscribe_types_id_seq',
        (
            SELECT
                max(id)
            FROM
                subscribe_types
        )
    );
