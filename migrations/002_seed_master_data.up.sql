-- Seed master_brands
INSERT INTO master_brands (name) VALUES
('TECKHING'),
('UNINEST'),
('TECHKING'),
('HILO'),
('BRIDGESTONE'),
('ADVANCE'),
('VK TYRE'),
('TRIANGLE'),
('Giti'),
('SAKURA');

-- Seed master_sizes
INSERT INTO master_sizes (name) VALUES
('16.00R25'),
('14.00R24'),
('20.5R25'),
('14.00-25'),
('14.00X24'),
('20.5X25');

-- Seed master_types
INSERT INTO master_types (name) VALUES
('Radial'),
('Bias');

-- Seed master_reasons
INSERT INTO master_reasons (name) VALUES
('Inner Liner Separation'),
('Bead Separation'),
('Bulging'),
('Tread Separation'),
('Rim Crack'),
('Tread Puncture'),
('Matching'),
('Bulging Chaffer'),
('Impact Damage'),
('Spud Grommets Broken'),
('Spud Grommets Loose'),
('Smooth Sidewall Cut'),
('Sidewall Damage'),
('Overheating'),
('Nail/Puncture'),
('Uneven Wear');

-- Seed master_actions
INSERT INTO master_actions (name) VALUES
('Mount'),
('Dismount'),
('Swap to Pos 1'),
('Swap to Pos 2'),
('Swap to Pos 3'),
('Swap to Pos 4'),
('Swap to Pos 5'),
('Swap to Pos 6'),
('Swap to Pos 7'),
('Swap to Pos 8'),
('Swap to Pos 9'),
('Swap to Pos 10'),
('Repair'),
('Scrap');

-- Seed master_remarks
INSERT INTO master_remarks (name) VALUES
('New Tyre'),
('Repair'),
('Running'),
('Spare'),
('Scrap'),
('UNIT SECOND HAND'),
('Ex Repair'),
('Canibal Tyre');

-- Seed unit_type_configs (SANY 10 positions)
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config) VALUES (
    'SANY_10POS',
    'SANY SKT 105S (10 Posisi)',
    10,
    '[{"position": 1, "label": "R-L1", "side": "rear_left", "axle": "rear_1", "x": 0.10, "y": 0.25},
      {"position": 2, "label": "R-L2", "side": "rear_left", "axle": "rear_1", "x": 0.25, "y": 0.25},
      {"position": 3, "label": "R-L3", "side": "rear_left", "axle": "rear_2", "x": 0.40, "y": 0.25},
      {"position": 4, "label": "R-R1", "side": "rear_right", "axle": "rear_1", "x": 0.60, "y": 0.25},
      {"position": 5, "label": "R-R2", "side": "rear_right", "axle": "rear_2", "x": 0.75, "y": 0.25},
      {"position": 6, "label": "R-R3", "side": "rear_right", "axle": "rear_2", "x": 0.90, "y": 0.25},
      {"position": 7, "label": "R-L4", "side": "rear_left", "axle": "bogie", "x": 0.17, "y": 0.55},
      {"position": 8, "label": "R-R4", "side": "rear_right", "axle": "bogie", "x": 0.83, "y": 0.55},
      {"position": 9, "label": "F-L1", "side": "front_left", "axle": "front", "x": 0.35, "y": 0.75},
      {"position": 10, "label": "F-R1", "side": "front_right", "axle": "front", "x": 0.65, "y": 0.75}]'::jsonb
);

-- Seed unit_type_configs (GREADER 6 positions)
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config) VALUES (
    'GREADER_6POS',
    'Greader (6 Posisi)',
    6,
    '[{"position": 1, "label": "R-L", "side": "rear_left", "axle": "rear", "x": 0.20, "y": 0.30},
      {"position": 2, "label": "R-R", "side": "rear_right", "axle": "rear", "x": 0.80, "y": 0.30},
      {"position": 3, "label": "F-L", "side": "front_left", "axle": "front", "x": 0.30, "y": 0.70},
      {"position": 4, "label": "F-R", "side": "front_right", "axle": "front", "x": 0.70, "y": 0.70},
      {"position": 5, "label": "R-L2", "side": "rear_left", "axle": "rear_2", "x": 0.20, "y": 0.55},
      {"position": 6, "label": "R-R2", "side": "rear_right", "axle": "rear_2", "x": 0.80, "y": 0.55}]'::jsonb
);

-- Seed unit_type_configs (ADT 8 positions)
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config) VALUES (
    'ADT_8POS',
    'Articulated Dump Truck (8 Posisi)',
    8,
    '[{"position": 1, "label": "R-L1", "side": "rear_left", "axle": "rear_1", "x": 0.15, "y": 0.25},
      {"position": 2, "label": "R-R1", "side": "rear_right", "axle": "rear_1", "x": 0.85, "y": 0.25},
      {"position": 3, "label": "R-L2", "side": "rear_left", "axle": "rear_2", "x": 0.25, "y": 0.25},
      {"position": 4, "label": "R-R2", "side": "rear_right", "axle": "rear_2", "x": 0.75, "y": 0.25},
      {"position": 5, "label": "R-L3", "side": "rear_left", "axle": "bogie", "x": 0.20, "y": 0.55},
      {"position": 6, "label": "R-R3", "side": "rear_right", "axle": "bogie", "x": 0.80, "y": 0.55},
      {"position": 7, "label": "F-L", "side": "front_left", "axle": "front", "x": 0.35, "y": 0.75},
      {"position": 8, "label": "F-R", "side": "front_right", "axle": "front", "x": 0.65, "y": 0.75}]'::jsonb
);

-- Seed master_patterns (after brands are inserted)
INSERT INTO master_patterns (brand_id, name)
SELECT b.id, p.name FROM master_brands b
CROSS JOIN (VALUES
    ('TECKHING', 'ET919'),
    ('TECKHING', 'ET919+'),
    ('TECKHING', 'ET668'),
    ('UNINEST', 'TIBERUN 811'),
    ('TECHKING', 'ET919'),
    ('TECHKING', 'ET919+'),
    ('TECHKING', 'VUT'),
    ('HILO', 'B01NL'),
    ('BRIDGESTONE', 'VUT'),
    ('ADVANCE', 'V-LUG'),
    ('VK TYRE', 'XTRA LOAD GRIP'),
    ('TRIANGLE', 'TB 516S'),
    ('Giti', 'GAO802')
) AS p(brand_name, name) WHERE b.name = p.brand_name;

-- Seed superadmin user (password: password123)
-- bcrypt hash for 'password123': $2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm
INSERT INTO users (name, email, password, role, company_id, status) VALUES
('Super Admin', 'admin@tms.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'superadmin', NULL, 'active');
