-- ============================================================
-- Migration: 001_init_schema
-- Consolidated: schema + master seed + company seed
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- COMPANIES
-- ============================================================
CREATE TABLE IF NOT EXISTS companies (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    address         TEXT NULL,
    contact_person  VARCHAR(255) NULL,
    phone           VARCHAR(50) NULL,
    email           VARCHAR(255) NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_companies_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_companies_name ON companies(name);
CREATE INDEX IF NOT EXISTS idx_companies_status ON companies(status);

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'admin_company',
    company_id      BIGINT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_users_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE SET NULL,
    CONSTRAINT chk_users_role CHECK (role IN ('superadmin', 'admin_company')),
    CONSTRAINT chk_users_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- ============================================================
-- PROJECTS
-- ============================================================
CREATE TABLE IF NOT EXISTS projects (
    id              BIGSERIAL PRIMARY KEY,
    company_id      BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    location        TEXT NULL,
    start_date      DATE NULL,
    end_date        DATE NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_projects_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT chk_projects_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_projects_company_id ON projects(company_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);

-- ============================================================
-- DRIVERS
-- ============================================================
CREATE TABLE IF NOT EXISTS drivers (
    id              BIGSERIAL PRIMARY KEY,
    company_id      BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    employee_id     VARCHAR(50) NOT NULL,
    phone           VARCHAR(50) NULL,
    license_number  VARCHAR(50) NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_drivers_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT chk_drivers_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_drivers_company_id ON drivers(company_id);
CREATE INDEX IF NOT EXISTS idx_drivers_employee_id ON drivers(employee_id);
CREATE INDEX IF NOT EXISTS idx_drivers_status ON drivers(status);

-- ============================================================
-- MASTER BRANDS
-- ============================================================
CREATE TABLE IF NOT EXISTS master_brands (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_brands_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- MASTER SIZES
-- ============================================================
CREATE TABLE IF NOT EXISTS master_sizes (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_sizes_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- MASTER TYPES
-- ============================================================
CREATE TABLE IF NOT EXISTS master_types (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_types_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- MASTER PATTERNS
-- ============================================================
CREATE TABLE IF NOT EXISTS master_patterns (
    id          BIGSERIAL PRIMARY KEY,
    brand_id    BIGINT NOT NULL,
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_mp_brand FOREIGN KEY (brand_id)
        REFERENCES master_brands(id) ON DELETE RESTRICT,
    CONSTRAINT uq_master_patterns_brand_name UNIQUE (brand_id, name),
    CONSTRAINT chk_master_patterns_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_master_patterns_brand_id ON master_patterns(brand_id);

-- ============================================================
-- MASTER REASONS
-- ============================================================
CREATE TABLE IF NOT EXISTS master_reasons (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_reasons_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- MASTER ACTIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS master_actions (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_actions_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- MASTER REMARKS
-- ============================================================
CREATE TABLE IF NOT EXISTS master_remarks (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL UNIQUE,
    code        VARCHAR(50) NULL,
    description TEXT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_master_remarks_status CHECK (status IN ('active', 'inactive'))
);

-- ============================================================
-- UNIT TYPE CONFIGS
-- ============================================================
CREATE TABLE IF NOT EXISTS unit_type_configs (
    id              BIGSERIAL PRIMARY KEY,
    unit_type       VARCHAR(50) NOT NULL UNIQUE,
    display_name    VARCHAR(100) NOT NULL,
    max_position    INTEGER NOT NULL,
    position_config JSONB NOT NULL,
    description     TEXT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================
-- UNITS
-- ============================================================
CREATE TABLE IF NOT EXISTS units (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    project_id          BIGINT NOT NULL,
    unit_id             VARCHAR(50) NOT NULL UNIQUE,
    unit_model          VARCHAR(255) NOT NULL,
    plate_number        VARCHAR(50) NULL,
    tyre_size_default   VARCHAR(50) NOT NULL,
    unit_type           VARCHAR(50) NOT NULL DEFAULT 'ADT',
    max_position        INTEGER NOT NULL DEFAULT 6,
    current_hm          DECIMAL(15,2) NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_units_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT fk_units_project FOREIGN KEY (project_id)
        REFERENCES projects(id) ON DELETE RESTRICT,
    CONSTRAINT chk_units_max_position CHECK (max_position > 0 AND max_position <= 20),
    CONSTRAINT chk_units_status CHECK (status IN ('active', 'inactive', 'maintenance'))
);
CREATE INDEX IF NOT EXISTS idx_units_company_id ON units(company_id);
CREATE INDEX IF NOT EXISTS idx_units_project_id ON units(project_id);
CREATE INDEX IF NOT EXISTS idx_units_unit_id ON units(unit_id);
CREATE INDEX IF NOT EXISTS idx_units_unit_type ON units(unit_type);
CREATE INDEX IF NOT EXISTS idx_units_status ON units(status);

-- ============================================================
-- TYRE MASTER
-- ============================================================
CREATE TABLE IF NOT EXISTS tyre_master (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    unit_id             BIGINT NULL,
    mounted_position    VARCHAR(20) NULL,
    barcode             VARCHAR(100) NOT NULL UNIQUE,
    serial_number       VARCHAR(100) NOT NULL UNIQUE,
    dot_code            VARCHAR(100) NULL,
    type                VARCHAR(50) NOT NULL DEFAULT 'Radial',
    size_id             BIGINT NOT NULL,
    brand_id            BIGINT NOT NULL,
    pattern_id          BIGINT NOT NULL,
    otd                 DECIMAL(5,2) NOT NULL DEFAULT 0,
    rtd                 DECIMAL(5,2) NOT NULL DEFAULT 0,
    rtd1                DECIMAL(5,2) NULL,
    rtd2                DECIMAL(5,2) NULL,
    lifetime            DECIMAL(15,2) NOT NULL DEFAULT 0,
    psi                 DECIMAL(5,2) NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'spare',
    remarks             TEXT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_tyre_master_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT fk_tyre_master_unit FOREIGN KEY (unit_id)
        REFERENCES units(id) ON DELETE SET NULL,
    CONSTRAINT fk_tyre_master_size FOREIGN KEY (size_id)
        REFERENCES master_sizes(id) ON DELETE RESTRICT,
    CONSTRAINT fk_tyre_master_brand FOREIGN KEY (brand_id)
        REFERENCES master_brands(id) ON DELETE RESTRICT,
    CONSTRAINT fk_tyre_master_pattern FOREIGN KEY (pattern_id)
        REFERENCES master_patterns(id) ON DELETE RESTRICT,
    CONSTRAINT chk_tyre_master_status CHECK (status IN ('spare', 'mounted', 'dismounted', 'scrap')),
    CONSTRAINT chk_tyre_master_rtd CHECK (rtd >= 0),
    CONSTRAINT chk_tyre_master_otd CHECK (otd >= 0),
    CONSTRAINT chk_tyre_master_psi CHECK (psi IS NULL OR psi > 0)
);
CREATE INDEX IF NOT EXISTS idx_tyre_master_company_id ON tyre_master(company_id);
CREATE INDEX IF NOT EXISTS idx_tyre_master_unit_id ON tyre_master(unit_id);
CREATE INDEX IF NOT EXISTS idx_tyre_master_barcode ON tyre_master(barcode);
CREATE INDEX IF NOT EXISTS idx_tyre_master_serial_number ON tyre_master(serial_number);
CREATE INDEX IF NOT EXISTS idx_tyre_master_status ON tyre_master(status);
CREATE INDEX IF NOT EXISTS idx_tyre_master_brand_id ON tyre_master(brand_id);
CREATE INDEX IF NOT EXISTS idx_tyre_master_size_id ON tyre_master(size_id);
CREATE INDEX IF NOT EXISTS idx_tyre_master_pattern_id ON tyre_master(pattern_id);

-- ============================================================
-- REPLACEMENTS
-- ============================================================
CREATE TABLE IF NOT EXISTS replacements (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    project_id          BIGINT NOT NULL,
    unit_id             BIGINT NOT NULL,
    driver_id           BIGINT NOT NULL,
    date                DATE NOT NULL,
    hm_update           DECIMAL(15,2) NOT NULL DEFAULT 0,
    current_life_hm     DECIMAL(15,2) NOT NULL,
    hm_plan             DECIMAL(15,2) NOT NULL,
    remarks             TEXT NULL,
    created_by          BIGINT NOT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_replacements_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT fk_replacements_project FOREIGN KEY (project_id)
        REFERENCES projects(id) ON DELETE RESTRICT,
    CONSTRAINT fk_replacements_unit FOREIGN KEY (unit_id)
        REFERENCES units(id) ON DELETE RESTRICT,
    CONSTRAINT fk_replacements_driver FOREIGN KEY (driver_id)
        REFERENCES drivers(id) ON DELETE RESTRICT,
    CONSTRAINT fk_replacements_created_by FOREIGN KEY (created_by)
        REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT chk_replacements_hm_plan CHECK (hm_plan > current_life_hm)
);
CREATE INDEX IF NOT EXISTS idx_replacements_company_id ON replacements(company_id);
CREATE INDEX IF NOT EXISTS idx_replacements_project_id ON replacements(project_id);
CREATE INDEX IF NOT EXISTS idx_replacements_unit_id ON replacements(unit_id);
CREATE INDEX IF NOT EXISTS idx_replacements_driver_id ON replacements(driver_id);
CREATE INDEX IF NOT EXISTS idx_replacements_date ON replacements(date);
CREATE INDEX IF NOT EXISTS idx_replacements_created_at ON replacements(created_at);

-- ============================================================
-- REPLACEMENT DETAILS
-- ============================================================
CREATE TABLE IF NOT EXISTS replacement_details (
    id                          BIGSERIAL PRIMARY KEY,
    replacement_id              BIGINT NOT NULL,
    position                    VARCHAR(20) NOT NULL,
    action                      VARCHAR(50) NOT NULL,
    old_tyre_id                BIGINT NULL,
    old_tyre_serial_number     VARCHAR(100) NULL,
    old_tyre_pattern           VARCHAR(100) NULL,
    old_tyre_size              VARCHAR(50) NULL,
    old_tyre_tread_1           DECIMAL(5,2) NULL,
    old_tyre_tread_2           DECIMAL(5,2) NULL,
    old_tyre_lifetime          DECIMAL(15,2) NULL,
    old_tyre_status            VARCHAR(20) NULL,
    failure_reason_id           BIGINT NULL,
    from_unit_id               VARCHAR(50) NULL,
    new_tyre_id                BIGINT NULL,
    new_tyre_serial_number     VARCHAR(100) NULL,
    new_tyre_pattern           VARCHAR(100) NULL,
    new_tyre_size              VARCHAR(50) NULL,
    new_tyre_tread_1           DECIMAL(5,2) NULL,
    new_tyre_tread_2           DECIMAL(5,2) NULL,
    new_tyre_current_lifetime  DECIMAL(15,2) NOT NULL DEFAULT 0,
    new_tyre_status            VARCHAR(50) NULL,
    remark                     TEXT NULL,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_rd_replacement FOREIGN KEY (replacement_id)
        REFERENCES replacements(id) ON DELETE CASCADE,
    CONSTRAINT fk_rd_old_tyre FOREIGN KEY (old_tyre_id)
        REFERENCES tyre_master(id) ON DELETE SET NULL,
    CONSTRAINT fk_rd_new_tyre FOREIGN KEY (new_tyre_id)
        REFERENCES tyre_master(id) ON DELETE SET NULL,
    CONSTRAINT fk_rd_failure_reason FOREIGN KEY (failure_reason_id)
        REFERENCES master_reasons(id) ON DELETE SET NULL,
    CONSTRAINT chk_rd_action CHECK (action IN ('mount', 'dismount', 'swap'))
);
CREATE INDEX IF NOT EXISTS idx_rd_replacement_id ON replacement_details(replacement_id);
CREATE INDEX IF NOT EXISTS idx_rd_old_tyre_id ON replacement_details(old_tyre_id);
CREATE INDEX IF NOT EXISTS idx_rd_new_tyre_id ON replacement_details(new_tyre_id);
CREATE INDEX IF NOT EXISTS idx_rd_failure_reason_id ON replacement_details(failure_reason_id);

-- ============================================================
-- MASTER DATA SEED
-- ============================================================

-- Seed master_brands
INSERT INTO master_brands (name) VALUES
('TECKHING'),('UNINEST'),('TECHKING'),('HILO'),('BRIDGESTONE'),('ADVANCE'),('VK TYRE'),('TRIANGLE'),('Giti'),('SAKURA')
ON CONFLICT (name) DO NOTHING;

-- Seed master_sizes
INSERT INTO master_sizes (name) VALUES
('16.00R25'),('14.00R24'),('20.5R25'),('14.00-25'),('14.00X24'),('20.5X25')
ON CONFLICT (name) DO NOTHING;

-- Seed master_types
INSERT INTO master_types (name) VALUES
('Radial'),('Bias')
ON CONFLICT (name) DO NOTHING;

-- Seed master_reasons
INSERT INTO master_reasons (name) VALUES
('Inner Liner Separation'),('Bead Separation'),('Bulging'),('Tread Separation'),('Rim Crack'),('Tread Puncture'),('Matching'),('Bulging Chaffer'),('Impact Damage'),('Spud Grommets Broken'),('Spud Grommets Loose'),('Smooth Sidewall Cut'),('Sidewall Damage'),('Overheating'),('Nail/Puncture'),('Uneven Wear')
ON CONFLICT (name) DO NOTHING;

-- Seed master_actions
INSERT INTO master_actions (name) VALUES
('Mount'),('Dismount'),('Swap to Pos 1'),('Swap to Pos 2'),('Swap to Pos 3'),('Swap to Pos 4'),('Swap to Pos 5'),('Swap to Pos 6'),('Swap to Pos 7'),('Swap to Pos 8'),('Swap to Pos 9'),('Swap to Pos 10'),('Repair'),('Scrap')
ON CONFLICT (name) DO NOTHING;

-- Seed master_remarks
INSERT INTO master_remarks (name) VALUES
('New Tyre'),('Repair'),('Running'),('Spare'),('Scrap'),('UNIT SECOND HAND'),('Ex Repair'),('Canibal Tyre')
ON CONFLICT (name) DO NOTHING;

-- ============================================================
-- UNIT TYPE CONFIGS SEED (final format: poros_1/2/3, Tyre N labels)
-- ============================================================

-- ADT_10POS: 8 positions across 3 axles
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config, status) VALUES (
    'ADT_10POS',
    'Articulated Dump Truck (8 Pos)',
    8,
    '[{"position":"1","label":"Tyre 1","side":"left","axle":"poros_3","x":0.15,"y":0.25,"mirror_of":"2"},
     {"position":"2","label":"Tyre 2","side":"right","axle":"poros_3","x":0.85,"y":0.25,"mirror_of":"1"},
     {"position":"3","label":"Tyre 3","side":"left","axle":"poros_4","x":0.25,"y":0.25,"mirror_of":"4"},
     {"position":"4","label":"Tyre 4","side":"right","axle":"poros_4","x":0.75,"y":0.25,"mirror_of":"3"},
     {"position":"5","label":"Tyre 5","side":"left","axle":"poros_2","x":0.20,"y":0.55,"mirror_of":"6"},
     {"position":"6","label":"Tyre 6","side":"right","axle":"poros_2","x":0.80,"y":0.55,"mirror_of":"5"},
     {"position":"7","label":"Tyre 7","side":"left","axle":"poros_1","x":0.35,"y":0.75,"mirror_of":"8"},
     {"position":"8","label":"Tyre 8","side":"right","axle":"poros_1","x":0.65,"y":0.75,"mirror_of":"7"}]'::jsonb,
    'active'
)
ON CONFLICT (unit_type) DO NOTHING;

-- DUMP_6POS: 6 positions across 2 axles (FIXED: max_position=6, all 6 entries)
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config, status) VALUES (
    'DUMP_6POS',
    'Dump Truck (6 Pos)',
    6,
    '[{"position":"1","label":"Tyre 1","side":"left","axle":"poros_1","x":0.2425,"y":0.1542,"mirror_of":"2"},
     {"position":"2","label":"Tyre 2","side":"right","axle":"poros_1","x":0.7575,"y":0.1542,"mirror_of":"1"},
     {"position":"3","label":"Tyre 3","side":"left","axle":"poros_2","x":0.1774,"y":0.6082,"mirror_of":"6"},
     {"position":"4","label":"Tyre 4","side":"left","axle":"poros_2","x":0.299,"y":0.6069,"mirror_of":"5"},
     {"position":"5","label":"Tyre 5","side":"right","axle":"poros_2","x":0.701,"y":0.6069,"mirror_of":"4"},
     {"position":"6","label":"Tyre 6","side":"right","axle":"poros_2","x":0.8226,"y":0.6082,"mirror_of":"3"}]'::jsonb,
    'active'
)
ON CONFLICT (unit_type) DO NOTHING;

-- COMPACT_4POS: 4 positions across 2 axles
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config, status) VALUES (
    'COMPACT_4POS',
    'Compact Vehicle (4 Pos)',
    4,
    '[{"position":"1","label":"Tyre 1","side":"left","axle":"poros_1","x":0.2863,"y":0.1569,"mirror_of":"2"},
     {"position":"2","label":"Tyre 2","side":"right","axle":"poros_1","x":0.7137,"y":0.1569,"mirror_of":"1"},
     {"position":"3","label":"Tyre 3","side":"left","axle":"poros_2","x":0.2987,"y":0.7972,"mirror_of":"4"},
     {"position":"4","label":"Tyre 4","side":"right","axle":"poros_2","x":0.7013,"y":0.7972,"mirror_of":"3"}]'::jsonb,
    'active'
)
ON CONFLICT (unit_type) DO NOTHING;

-- TRUCK_6POS: 6 positions across 3 axles
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config, status) VALUES (
    'TRUCK_6POS',
    'Standard Truck (6 Pos)',
    6,
    '[{"position":"1","label":"Tyre 1","side":"left","axle":"poros_1","x":0.2654,"y":0.1556,"mirror_of":"2"},
     {"position":"2","label":"Tyre 2","side":"right","axle":"poros_1","x":0.7346,"y":0.1556,"mirror_of":"1"},
     {"position":"3","label":"Tyre 3","side":"left","axle":"poros_2","x":0.2675,"y":0.6056,"mirror_of":"4"},
     {"position":"4","label":"Tyre 4","side":"right","axle":"poros_2","x":0.7325,"y":0.6056,"mirror_of":"3"},
     {"position":"5","label":"Tyre 5","side":"left","axle":"poros_3","x":0.2675,"y":0.7914,"mirror_of":"6"},
     {"position":"6","label":"Tyre 6","side":"right","axle":"poros_3","x":0.7325,"y":0.7914,"mirror_of":"5"}]'::jsonb,
    'active'
)
ON CONFLICT (unit_type) DO NOTHING;

-- Seed master_patterns
INSERT INTO master_patterns (brand_id, name)
SELECT b.id, p.name FROM master_brands b
CROSS JOIN (VALUES
    ('TECKHING', 'ET919'),('TECKHING', 'ET919+'),('TECKHING', 'ET668'),
    ('UNINEST', 'TIBERUN 811'),
    ('TECHKING', 'ET919'),('TECHKING', 'ET919+'),('TECHKING', 'VUT'),
    ('HILO', 'B01NL'),
    ('BRIDGESTONE', 'VUT'),
    ('ADVANCE', 'V-LUG'),
    ('VK TYRE', 'XTRA LOAD GRIP'),
    ('TRIANGLE', 'TB 516S'),
    ('Giti', 'GAO802')
) AS p(brand_name, name) WHERE b.name = p.brand_name
ON CONFLICT ON CONSTRAINT uq_master_patterns_brand_name DO NOTHING;

-- Seed superadmin user (password: password123)
INSERT INTO users (name, email, password, role, company_id, status) VALUES
('Super Admin', 'admin@tms.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'superadmin', NULL, 'active')
ON CONFLICT (email) DO NOTHING;

-- ============================================================
-- COMPANY DATA SEED
-- ============================================================

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

-- Seed Units (3x SANY SKT 105S, ADT_10POS, max_position=8)
INSERT INTO units (company_id, project_id, unit_id, unit_model, plate_number, tyre_size_default, unit_type, max_position, current_hm, status)
SELECT c.id, p.id, u.unit_id, u.unit_model, u.plate, u.size, 'ADT_10POS', 8, u.hm, 'active'
FROM companies c
CROSS JOIN projects p
CROSS JOIN (VALUES
    ('BWB001', 'SANY SKT 105S', 'KT 1234 AB', '16.00R25', 12850.00),
    ('BWB002', 'SANY SKT 105S', 'KT 5678 CD', '16.00R25', 11420.50),
    ('BWB003', 'SANY SKT 105S', 'KT 9012 EF', '16.00R25', 9875.25)
) AS u(unit_id, unit_model, plate, size, hm)
WHERE c.name = 'Berkat Anuegrah Sejahtera' AND p.name = 'BSSR'
ON CONFLICT DO NOTHING;

-- Seed Tyres (15 mounted + 5 spare, numeric position key matching position_config)
-- Mounted tyres: join company + unit by unit_code
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
    t.barcode, t.serial, t.dot,
    t.tyre_type, sz.id, b.id, pt.id,
    t.otd, t.rtd, t.rtd1, t.rtd2, t.lifetime, t.psi,
    t.status, t.remarks
FROM companies c
CROSS JOIN units un
CROSS JOIN (VALUES
    ('BWB001','1','TYR00001','SN00001','DOT-2022-001','Radial','16.00R25','TECKHING','ET919',28.5,24.0,23.5,24.5,3250.0,95.0,'mounted','Front left - rear axle 1'),
    ('BWB001','2','TYR00002','SN00002','DOT-2022-002','Radial','16.00R25','TECHKING','ET919',26.0,22.0,21.0,23.0,4100.5,90.0,'mounted','Front left - rear axle 1'),
    ('BWB001','3','TYR00003','SN00003','DOT-2022-003','Radial','16.00R25','TECKHING','ET919+',32.0,28.5,29.0,28.0,1800.0,95.0,'mounted','Front left - rear axle 2'),
    ('BWB001','4','TYR00004','SN00004','DOT-2022-004','Radial','16.00R25','HILO','B01NL',18.0,14.5,13.0,16.0,7200.0,85.0,'mounted','Front right - rear axle 1'),
    ('BWB001','5','TYR00005','SN00005','DOT-2022-005','Radial','16.00R25','BRIDGESTONE','VUT',15.5,11.0,10.5,11.5,9500.0,80.0,'mounted','Front right - rear axle 2'),
    ('BWB002','1','TYR00011','SN00011','DOT-2022-011','Radial','16.00R25','TECHKING','ET919',25.0,21.5,20.0,23.0,4350.0,90.0,'mounted','Front left - rear axle 1'),
    ('BWB002','2','TYR00012','SN00012','DOT-2022-012','Radial','16.00R25','TECKHING','ET919+',30.0,26.0,25.5,26.5,2100.0,95.0,'mounted','Front left - rear axle 1'),
    ('BWB002','3','TYR00013','SN00013','DOT-2022-013','Radial','16.00R25','TRIANGLE','TB 516S',22.0,18.0,17.0,19.0,5800.0,88.0,'mounted','Front left - rear axle 2'),
    ('BWB002','4','TYR00014','SN00014','DOT-2022-014','Radial','16.00R25','HILO','B01NL',17.5,13.5,12.0,15.0,8100.0,85.0,'mounted','Front right - rear axle 1'),
    ('BWB002','5','TYR00015','SN00015','DOT-2022-015','Radial','16.00R25','Giti','GAO802',20.0,16.5,15.0,18.0,6500.0,88.0,'mounted','Front right - rear axle 2'),
    ('BWB003','1','TYR00016','SN00016','DOT-2022-016','Radial','16.00R25','TECKHING','ET919',27.5,23.0,22.5,23.5,3800.0,95.0,'mounted','Front left - rear axle 1'),
    ('BWB003','2','TYR00017','SN00017','DOT-2022-017','Radial','16.00R25','TECKHING','ET919+',33.0,29.5,30.0,29.0,1500.0,95.0,'mounted','Front left - rear axle 1'),
    ('BWB003','3','TYR00018','SN00018','DOT-2022-018','Radial','16.00R25','BRIDGESTONE','VUT',19.0,15.0,14.0,16.0,7500.0,82.0,'mounted','Front left - rear axle 2'),
    ('BWB003','4','TYR00019','SN00019','DOT-2022-019','Radial','16.00R25','TECHKING','VUT',16.0,12.0,11.0,13.0,9200.0,80.0,'mounted','Front right - rear axle 1'),
    ('BWB003','5','TYR00020','SN00020','DOT-2022-020','Radial','16.00R25','UNINEST','TIBERUN 811',24.0,20.0,19.0,21.0,5100.0,90.0,'mounted','Front right - rear axle 2')
) AS t(unit_code,mounted_pos,barcode,serial,dot,tyre_type,size_name,brand_name,pattern_name,otd,rtd,rtd1,rtd2,lifetime,psi,status,remarks)
LEFT JOIN master_sizes sz ON sz.name = t.size_name
LEFT JOIN master_brands b ON b.name = t.brand_name
LEFT JOIN master_patterns pt ON pt.name = t.pattern_name AND pt.brand_id = b.id
WHERE c.name = 'Berkat Anuegrah Sejahtera'
  AND un.unit_id = t.unit_code
ON CONFLICT DO NOTHING;

-- Spare tyres: no unit_id
INSERT INTO tyre_master (
    company_id, unit_id, mounted_position,
    barcode, serial_number, dot_code,
    type, size_id, brand_id, pattern_id,
    otd, rtd, rtd1, rtd2, lifetime, psi,
    status, remarks
)
SELECT
    c.id,
    NULL,
    t.mounted_pos,
    t.barcode, t.serial, t.dot,
    t.tyre_type, sz.id, b.id, pt.id,
    t.otd, t.rtd, t.rtd1, t.rtd2, t.lifetime, t.psi,
    t.status, t.remarks
FROM companies c
CROSS JOIN (VALUES
    ('Spare','TYR00006','SN00006','DOT-2022-006','Radial','16.00R25','TECKHING','ET668',38.0,36.0,36.5,35.5,500.0,98.0,'spare','New tyre - ready to mount'),
    ('Spare','TYR00007','SN00007','DOT-2022-007','Radial','16.00R25','UNINEST','TIBERUN 811',40.0,38.5,39.0,38.0,0.0,100.0,'spare','New tyre - ready to mount'),
    ('Spare','TYR00008','SN00008','DOT-2022-008','Radial','16.00R25','TECHKING','ET919+',39.0,37.5,38.0,37.0,200.0,98.0,'spare','New tyre - ready to mount'),
    ('Spare','TYR00009','SN00009','DOT-2022-009','Radial','16.00R25','TECKHING','ET919',35.0,33.0,32.5,33.5,850.0,95.0,'spare','Running low - monitor closely'),
    ('Spare','TYR00010','SN00010','DOT-2022-010','Radial','16.00R25','ADVANCE','V-LUG',30.0,27.0,26.5,27.5,1200.0,92.0,'spare','Ex repair - still usable')
) AS t(mounted_pos,barcode,serial,dot,tyre_type,size_name,brand_name,pattern_name,otd,rtd,rtd1,rtd2,lifetime,psi,status,remarks)
LEFT JOIN master_sizes sz ON sz.name = t.size_name
LEFT JOIN master_brands b ON b.name = t.brand_name
LEFT JOIN master_patterns pt ON pt.name = t.pattern_name AND pt.brand_id = b.id
WHERE c.name = 'Berkat Anuegrah Sejahtera'
ON CONFLICT DO NOTHING;
