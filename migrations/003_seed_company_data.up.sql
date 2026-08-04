-- Seed Company
INSERT INTO companies (name, address, contact_person, phone, email, status) VALUES
('Berkat Anuegrah Sejahtera', 'Kutai Kartanegara, Kalimantan Timur', 'Fikri Zufri', '0812-3456-7890', 'company@bas.com', 'active')
ON CONFLICT DO NOTHING;

-- Seed Project
INSERT INTO projects (company_id, name, location, start_date, end_date, status)
SELECT c.id, 'BSSR', 'Batuah, Kutai Kartanegara', '2024-01-01', '2026-12-31', 'active'
FROM companies c WHERE c.name = 'Berkat Anuegrah Sejahtera'
ON CONFLICT DO NOTHING;

-- Seed Drivers
INSERT INTO drivers (company_id, name, employee_id, phone, license_number, status)
SELECT c.id, d.name, d.emp_id, d.phone, d.license, 'active'
FROM companies c
CROSS JOIN (VALUES
    ('Rudi Hartono', 'DRV001', '0813-1111-2222', 'SIM B-1234-KT'),
    ('Surya Darma', 'DRV002', '0813-3333-4444', 'SIM B-5678-KT'),
    ('Asep Saepulloh', 'DRV003', '0813-5555-6666', 'SIM B-9012-KT')
) AS d(name, emp_id, phone, license)
WHERE c.name = 'Berkat Anuegrah Sejahtera'
ON CONFLICT DO NOTHING;

-- Seed Units (3x SANY SKT 105S, 10 positions each)
INSERT INTO units (company_id, project_id, unit_id, unit_model, plate_number, tyre_size_default, unit_type, max_position, current_hm, status)
SELECT c.id, p.id, u.unit_id, u.unit_model, u.plate, u.size, 'SANY_10POS', 10, u.hm, 'active'
FROM companies c
CROSS JOIN projects p
CROSS JOIN (VALUES
    ('BWB001', 'SANY SKT 105S', 'KT 1234 AB', '16.00R25', 12850.00),
    ('BWB002', 'SANY SKT 105S', 'KT 5678 CD', '16.00R25', 11420.50),
    ('BWB003', 'SANY SKT 105S', 'KT 9012 EF', '16.00R25', 9875.25)
) AS u(unit_id, unit_model, plate, size, hm)
WHERE c.name = 'Berkat Anuegrah Sejahtera'
  AND p.name = 'BSSR'
ON CONFLICT DO NOTHING;

-- Seed Tyres (20 tyres: 15 mounted across 3 units, 5 spare)
INSERT INTO tyre_master (
    company_id, unit_id, mounted_position,
    barcode, serial_number, dot_code,
    type, size_id, brand_id, pattern_id,
    otd, rtd, rtd1, rtd2, lifetime, psi,
    status, remarks
)
SELECT
    c.id,
    un.id,
    t.mounted_pos,
    t.barcode,
    t.serial,
    t.dot,
    t.tyre_type,
    sz.id,
    b.id,
    pt.id,
    t.otd,
    t.rtd,
    t.rtd1,
    t.rtd2,
    t.lifetime,
    t.psi,
    t.status,
    t.remarks
FROM companies c
CROSS JOIN units un
CROSS JOIN (VALUES
    -- Unit BWB001 tyres (5 mounted + 5 spare)
    ('BWB001', 'Pos 1',  'TYR00001', 'SN00001', 'DOT-2022-001', 'Radial', '16.00R25', 'TECKHING',  'ET919',     28.5, 24.0, 23.5, 24.5, 3250.0, 95.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB001', 'Pos 2',  'TYR00002', 'SN00002', 'DOT-2022-002', 'Radial', '16.00R25', 'TECHKING',  'ET919',     26.0, 22.0, 21.0, 23.0, 4100.5, 90.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB001', 'Pos 3',  'TYR00003', 'SN00003', 'DOT-2022-003', 'Radial', '16.00R25', 'TECKHING',  'ET919+',    32.0, 28.5, 29.0, 28.0, 1800.0, 95.0, 'mounted', 'Front left — rear axle 2'),
    ('BWB001', 'Pos 4',  'TYR00004', 'SN00004', 'DOT-2022-004', 'Radial', '16.00R25', 'HILO',      'B01NL',     18.0, 14.5, 13.0, 16.0, 7200.0, 85.0, 'mounted', 'Front right — rear axle 1'),
    ('BWB001', 'Pos 5',  'TYR00005', 'SN00005', 'DOT-2022-005', 'Radial', '16.00R25', 'BRIDGESTONE', 'VUT',      15.5, 11.0, 10.5, 11.5, 9500.0, 80.0, 'mounted', 'Front right — rear axle 2'),
    (NULL,    'Spare',  'TYR00006', 'SN00006', 'DOT-2022-006', 'Radial', '16.00R25', 'TECKHING',  'ET668',     38.0, 36.0, 36.5, 35.5, 500.0,  98.0, 'spare',   'New tyre — ready to mount'),
    (NULL,    'Spare',  'TYR00007', 'SN00007', 'DOT-2022-007', 'Radial', '16.00R25', 'UNINEST',   'TIBERUN 811', 40.0, 38.5, 39.0, 38.0, 0.0, 100.0, 'spare', 'New tyre — ready to mount'),
    (NULL,    'Spare',  'TYR00008', 'SN00008', 'DOT-2022-008', 'Radial', '16.00R25', 'TECHKING',  'ET919+',    39.0, 37.5, 38.0, 37.0, 200.0, 98.0, 'spare',   'New tyre — ready to mount'),
    (NULL,    'Spare',  'TYR00009', 'SN00009', 'DOT-2022-009', 'Radial', '16.00R25', 'TECKHING',  'ET919',     35.0, 33.0, 32.5, 33.5, 850.0, 95.0, 'spare',   'Running low — monitor closely'),
    (NULL,    'Spare',  'TYR00010', 'SN00010', 'DOT-2022-010', 'Radial', '16.00R25', 'ADVANCE',   'V-LUG',     30.0, 27.0, 26.5, 27.5, 1200.0, 92.0, 'spare',   'Ex repair — still usable'),

    -- Unit BWB002 tyres (5 mounted)
    ('BWB002', 'Pos 1',  'TYR00011', 'SN00011', 'DOT-2022-011', 'Radial', '16.00R25', 'TECHKING',  'ET919',     25.0, 21.5, 20.0, 23.0, 4350.0, 90.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB002', 'Pos 2',  'TYR00012', 'SN00012', 'DOT-2022-012', 'Radial', '16.00R25', 'TECKHING',  'ET919+',    30.0, 26.0, 25.5, 26.5, 2100.0, 95.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB002', 'Pos 3',  'TYR00013', 'SN00013', 'DOT-2022-013', 'Radial', '16.00R25', 'TRIANGLE',  'TB 516S',   22.0, 18.0, 17.0, 19.0, 5800.0, 88.0, 'mounted', 'Front left — rear axle 2'),
    ('BWB002', 'Pos 4',  'TYR00014', 'SN00014', 'DOT-2022-014', 'Radial', '16.00R25', 'HILO',      'B01NL',     17.5, 13.5, 12.0, 15.0, 8100.0, 85.0, 'mounted', 'Front right — rear axle 1'),
    ('BWB002', 'Pos 5',  'TYR00015', 'SN00015', 'DOT-2022-015', 'Radial', '16.00R25', 'Giti',      'GAO802',    20.0, 16.5, 15.0, 18.0, 6500.0, 88.0, 'mounted', 'Front right — rear axle 2'),

    -- Unit BWB003 tyres (5 mounted)
    ('BWB003', 'Pos 1',  'TYR00016', 'SN00016', 'DOT-2022-016', 'Radial', '16.00R25', 'TECKHING',  'ET919',     27.5, 23.0, 22.5, 23.5, 3800.0, 95.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB003', 'Pos 2',  'TYR00017', 'SN00017', 'DOT-2022-017', 'Radial', '16.00R25', 'TECKHING',  'ET919+',    33.0, 29.5, 30.0, 29.0, 1500.0, 95.0, 'mounted', 'Front left — rear axle 1'),
    ('BWB003', 'Pos 3',  'TYR00018', 'SN00018', 'DOT-2022-018', 'Radial', '16.00R25', 'BRIDGESTONE', 'VUT',     19.0, 15.0, 14.0, 16.0, 7500.0, 82.0, 'mounted', 'Front left — rear axle 2'),
    ('BWB003', 'Pos 4',  'TYR00019', 'SN00019', 'DOT-2022-019', 'Radial', '16.00R25', 'TECHKING',  'VUT',       16.0, 12.0, 11.0, 13.0, 9200.0, 80.0, 'mounted', 'Front right — rear axle 1'),
    ('BWB003', 'Pos 5',  'TYR00020', 'SN00020', 'DOT-2022-020', 'Radial', '16.00R25', 'UNINEST',   'TIBERUN 811', 24.0, 20.0, 19.0, 21.0, 5100.0, 90.0, 'mounted', 'Front right — rear axle 2')
) AS t(unit_code, mounted_pos, barcode, serial, dot, tyre_type, size_name, brand_name, pattern_name, otd, rtd, rtd1, rtd2, lifetime, psi, status, remarks)
LEFT JOIN master_sizes  sz ON sz.name = t.size_name
LEFT JOIN master_brands b  ON b.name  = t.brand_name
LEFT JOIN master_patterns pt ON pt.name = t.pattern_name AND pt.brand_id = b.id
WHERE c.name = 'Berkat Anuegrah Sejahtera'
  AND un.unit_id = COALESCE(t.unit_code, '')
ON CONFLICT DO NOTHING;
