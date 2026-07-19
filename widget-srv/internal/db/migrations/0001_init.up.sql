CREATE TABLE product (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    attributes  JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory (
    product_id  UUID PRIMARY KEY REFERENCES product(id) ON DELETE CASCADE,
    count       INTEGER NOT NULL DEFAULT 0 CHECK (count >= 0)
);
