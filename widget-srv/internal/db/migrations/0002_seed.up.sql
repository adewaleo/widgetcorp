INSERT INTO product (id, name, attributes) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Standard Widget', '{"color": "blue", "size": "medium"}'),
    ('22222222-2222-2222-2222-222222222222', 'Deluxe Widget',   '{"color": "gold", "size": "large", "premium": true}'),
    ('33333333-3333-3333-3333-333333333333', 'Mini Widget',     '{"color": "red", "size": "small"}');

INSERT INTO inventory (product_id, count) VALUES
    ('11111111-1111-1111-1111-111111111111', 100),
    ('22222222-2222-2222-2222-222222222222', 25),
    ('33333333-3333-3333-3333-333333333333', 500);
