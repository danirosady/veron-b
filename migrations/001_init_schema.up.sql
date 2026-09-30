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
    name            VARCHAR(255) NOT NULL UNIQUE,
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

-- Add unique index for companies ON CONFLICT support (idempotent, skips if already exists)
CREATE UNIQUE INDEX IF NOT EXISTS uq_companies_name ON companies(name);

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
    CONSTRAINT uq_projects_company_name UNIQUE (company_id, name),
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
    CONSTRAINT uq_drivers_company_employee UNIQUE (company_id, employee_id),
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
('UNINEST'),('TECHKING'),('HILO'),('BRIDGESTONE'),('ADVANCE'),('VK TYRE'),('TRIANGLE'),('Giti'),('SAKURA')
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

-- WDT_10POS: 10 positions across 3 axles
INSERT INTO unit_type_configs (unit_type, display_name, max_position, position_config, status) VALUES (
    'WDT_10POS',
    'Wide Dump Truck (10 Pos)',
    10,
    '[{"position":"1","label":"Tyre 1","side":"left","axle":"poros_1","x":0.255,"y":0.1556,"mirror_of":"2"},
     {"position":"2","label":"Tyre 2","side":"right","axle":"poros_1","x":0.745,"y":0.1556,"mirror_of":"1"},
     {"position":"3","label":"Tyre 3","side":"left","axle":"poros_2","x":0.16,"y":0.61,"mirror_of":"6"},
     {"position":"4","label":"Tyre 4","side":"left","axle":"poros_2","x":0.2804,"y":0.6097,"mirror_of":"5"},
     {"position":"5","label":"Tyre 5","side":"right","axle":"poros_2","x":0.7196,"y":0.6097,"mirror_of":"4"},
     {"position":"6","label":"Tyre 6","side":"right","axle":"poros_2","x":0.84,"y":0.61,"mirror_of":"3"},
     {"position":"7","label":"Tyre 7","side":"left","axle":"poros_3","x":0.16,"y":0.8,"mirror_of":"10"},
     {"position":"8","label":"Tyre 8","side":"left","axle":"poros_3","x":0.28,"y":0.8,"mirror_of":"9"},
     {"position":"9","label":"Tyre 9","side":"right","axle":"poros_3","x":0.72,"y":0.8,"mirror_of":"8"},
     {"position":"10","label":"Tyre 10","side":"right","axle":"poros_3","x":0.84,"y":0.8,"mirror_of":"7"}]'::jsonb,
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
    ('UNINEST', 'TIBERUN 811'),
    ('TECHKING', 'ET919'),('TECHKING', 'ET919+'),('TECHKING', 'VUT'),
    ('HILO', 'B01NL'),
    ('BRIDGESTONE', 'VUT'),
    ('ADVANCE', 'V-LUG'),
    ('VK TYRE', 'XTRA LOAD GRIP'),
    ('TRIANGLE', 'TB 516S'),
    ('Giti', 'GAO802'),
    ('SAKURA', 'SRS-01')
) AS p(brand_name, name) WHERE b.name = p.brand_name
ON CONFLICT ON CONSTRAINT uq_master_patterns_brand_name DO NOTHING;

-- Seed superadmin user (password: password123)
INSERT INTO users (name, email, password, role, company_id, status) VALUES
('Super Admin', 'admin@tms.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'superadmin', NULL, 'active')
ON CONFLICT (email) DO NOTHING;

-- ============================================================
-- COMPANY DATA SEED (3 companies, 100 tyres total)
-- ============================================================

-- Seed Companies
INSERT INTO companies (name, address, contact_person, phone, email, status) VALUES
('Berkat Anuegrah Sejahtera', 'Kutai Kartanegara, Kalimantan Timur', 'Fikri Zufri', '0812-3456-7890', 'company@bas.com', 'active'),
('Borneo Energy Indonesia', 'Samarinda, Kalimantan Timur', 'Ahmad Fauzi', '0813-9876-5432', 'company@bei.com', 'active'),
('Kalimantan Prima Persada', 'Sangkulirang, Kutai Timur', 'Budi Santoso', '0815-1122-3344', 'company@kpp.com', 'active')
ON CONFLICT (name) DO NOTHING;

-- Seed Projects (2 per company)
INSERT INTO projects (company_id, name, location, start_date, end_date, status)
SELECT c.id, p.name, p.location, p.start_date, p.end_date, 'active'
FROM companies c
CROSS JOIN (VALUES
    -- BAS projects
    ('BSSR', 'Batuah, Kutai Kartanegara', DATE '2024-01-01', DATE '2026-12-31'),
    ('BSSM', 'Bontang, Kutai Kartanegara', DATE '2025-01-01', DATE '2027-06-30'),
    -- BEI projects
    ('EBI', 'Samarinda Utara, Samarinda', DATE '2024-06-01', DATE '2026-12-31'),
    ('EBS', 'Sengkotek, Samarinda', DATE '2025-03-01', DATE '2027-03-31'),
    -- KPP projects
    ('KPP', 'Sangkulirang, Kutai Timur', DATE '2024-09-01', DATE '2027-08-31'),
    ('KPS', 'Pantai Tanjung Bara, Kutai Timur', DATE '2025-06-01', DATE '2028-05-31')
) AS p(name, location, start_date, end_date)
WHERE (c.name = 'Berkat Anuegrah Sejahtera' AND p.name IN ('BSSR','BSSM'))
   OR (c.name = 'Borneo Energy Indonesia' AND p.name IN ('EBI','EBS'))
   OR (c.name = 'Kalimantan Prima Persada' AND p.name IN ('KPP','KPS'))
ON CONFLICT (company_id, name) DO NOTHING;

-- Seed admin_company users (1 per company, superadmin already exists)
INSERT INTO users (name, email, password, role, company_id, status) VALUES
('Admin BAS', 'admin@bas.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'admin_company', 1, 'active'),
('Admin BEI', 'admin@bei.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'admin_company', 2, 'active'),
('Admin KPP', 'admin@kpp.com', '$2a$10$Em30c27ErDXVIWBY0jopT.IsQTRYS4Kpd.Y792p.i1dVPOIgxootm', 'admin_company', 3, 'active')
ON CONFLICT (email) DO NOTHING;

-- Seed Drivers (3 per company)
INSERT INTO drivers (company_id, name, employee_id, phone, license_number, status)
SELECT c.id, d.name, d.emp_id, d.phone, d.license, 'active'
FROM companies c
CROSS JOIN (VALUES
    -- BAS drivers
    ('Rudi Hartono', 'DRV001', '0813-1111-2222', 'SIM B-1234-KT'),
    ('Surya Darma', 'DRV002', '0813-3333-4444', 'SIM B-5678-KT'),
    ('Asep Saepulloh', 'DRV003', '0813-5555-6666', 'SIM B-9012-KT'),
    -- BEI drivers
    ('Joko Widodo', 'DRV010', '0813-2222-3333', 'SIM B-2234-SM'),
    ('Dedi Kusuma', 'DRV011', '0813-4444-5555', 'SIM B-6678-SM'),
    ('Fajar Nugroho', 'DRV012', '0813-6666-7777', 'SIM B-0034-SM'),
    -- KPP drivers
    ('Heri Susanto', 'DRV020', '0813-7777-8888', 'SIM B-4456-KT'),
    ('Agus Salim', 'DRV021', '0813-8888-9999', 'SIM B-7789-KT'),
    ('Rizki Pratama', 'DRV022', '0813-9999-0000', 'SIM B-0090-KT')
) AS d(name, emp_id, phone, license)
WHERE (c.name = 'Berkat Anuegrah Sejahtera' AND d.emp_id IN ('DRV001','DRV002','DRV003'))
   OR (c.name = 'Borneo Energy Indonesia' AND d.emp_id IN ('DRV010','DRV011','DRV012'))
   OR (c.name = 'Kalimantan Prima Persada' AND d.emp_id IN ('DRV020','DRV021','DRV022'))
ON CONFLICT (company_id, employee_id) DO NOTHING;

-- Seed Units (3 per company, all SANY WDT_10POS, max_position=10)
INSERT INTO units (company_id, project_id, unit_id, unit_model, plate_number, tyre_size_default, unit_type, max_position, current_hm, status)
SELECT c.id, p.id, u.unit_id, u.unit_model, u.plate, u.size, 'WDT_10POS', 10, u.hm, 'active'
FROM companies c
JOIN projects p ON p.company_id = c.id
CROSS JOIN (VALUES
    -- BAS units
    ('BSSR', 'BWB001', 'SANY SKT 105S', 'KT 1234 AB', '16.00R25', 12850.00),
    ('BSSR', 'BWB002', 'SANY SKT 105S', 'KT 5678 CD', '16.00R25', 11420.50),
    ('BSSR', 'BWB003', 'SANY SKT 105S', 'KT 9012 EF', '16.00R25', 9875.25),
    -- BEI units
    ('EBI', 'BEI001', 'SANY SKT 105S', 'KT 2345 GH', '16.00R25', 15200.00),
    ('EBI', 'BEI002', 'SANY SKT 105S', 'KT 6789 IJ', '16.00R25', 13850.75),
    ('EBI', 'BEI003', 'SANY SKT 105S', 'KT 3456 KL', '16.00R25', 11200.00),
    -- KPP units
    ('KPP', 'KPP001', 'SANY SKT 105S', 'KT 7890 MN', '16.00R25', 16500.00),
    ('KPP', 'KPP002', 'SANY SKT 105S', 'KT 0123 OP', '16.00R25', 14100.25),
    ('KPP', 'KPP003', 'SANY SKT 105S', 'KT 4567 QR', '16.00R25', 10500.50)
) AS u(proj_name, unit_id, unit_model, plate, size, hm)
WHERE (c.name = 'Berkat Anuegrah Sejahtera' AND p.name = u.proj_name)
   OR (c.name = 'Borneo Energy Indonesia' AND p.name = u.proj_name)
   OR (c.name = 'Kalimantan Prima Persada' AND p.name = u.proj_name)
ON CONFLICT DO NOTHING;

-- Seed Tyres: 100 total across 3 companies
-- BAS: 30 tyres (Unit1=8 mounted, Unit2=4 mounted+4 spare, Unit3=14 spare)
-- BEI: 35 tyres (Unit1=8 mounted, Unit2=4 mounted+4 spare, Unit3=23 spare)
-- KPP: 35 tyres (Unit1=8 mounted, Unit2=4 mounted+4 spare, Unit3=23 spare)
-- Brands used: UNINEST, TECHKING, HILO, BRIDGESTONE, ADVANCE, VK TYRE, TRIANGLE, Giti, SAKURA
INSERT INTO tyre_master (
    company_id, unit_id, mounted_position,
    barcode, serial_number, dot_code,
    type, size_id, brand_id, pattern_id,
    otd, rtd, rtd1, rtd2, lifetime, psi,
    status, remarks
)
SELECT
    c.id,
    CASE WHEN t.unit_code = 'BWB001' THEN (SELECT id FROM units WHERE unit_id='BWB001' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'BWB002' THEN (SELECT id FROM units WHERE unit_id='BWB002' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'BWB003' THEN (SELECT id FROM units WHERE unit_id='BWB003' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'BEI001' THEN (SELECT id FROM units WHERE unit_id='BEI001' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'BEI002' THEN (SELECT id FROM units WHERE unit_id='BEI002' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'BEI003' THEN (SELECT id FROM units WHERE unit_id='BEI003' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'KPP001' THEN (SELECT id FROM units WHERE unit_id='KPP001' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'KPP002' THEN (SELECT id FROM units WHERE unit_id='KPP002' AND company_id=c.id LIMIT 1)
         WHEN t.unit_code = 'KPP003' THEN (SELECT id FROM units WHERE unit_id='KPP003' AND company_id=c.id LIMIT 1)
         ELSE NULL END,
    t.mounted_pos,
    t.barcode, t.serial, t.dot,
    t.tyre_type, sz.id, b.id, pt.id,
    t.otd, t.rtd, t.rtd1, t.rtd2, t.lifetime, t.psi,
    t.status, t.remarks
FROM companies c
CROSS JOIN (VALUES
    -- ========== BAS TYRES (30 total) ==========
    -- Unit 1 (BWB001): 8 mounted tyres - all positions
    ('BWB001','1','TYR00001','SN00001','DOT-2022-001','Radial','16.00R25','UNINEST','TIBERUN 811',30.0,26.0,25.5,26.5,1500.0,95.0,'mounted','Unit 1 - Pos 1'),
    ('BWB001','2','TYR00002','SN00002','DOT-2022-002','Radial','16.00R25','TECHKING','ET919',28.5,24.0,23.5,24.5,2200.0,90.0,'mounted','Unit 1 - Pos 2'),
    ('BWB001','3','TYR00003','SN00003','DOT-2022-003','Radial','16.00R25','HILO','B01NL',32.0,28.5,29.0,28.0,950.0,95.0,'mounted','Unit 1 - Pos 3'),
    ('BWB001','4','TYR00004','SN00004','DOT-2022-004','Radial','16.00R25','BRIDGESTONE','VUT',18.0,14.5,13.0,16.0,6800.0,85.0,'mounted','Unit 1 - Pos 4'),
    ('BWB001','5','TYR00005','SN00005','DOT-2022-005','Radial','16.00R25','ADVANCE','V-LUG',15.5,11.0,10.5,11.5,9100.0,80.0,'mounted','Unit 1 - Pos 5'),
    ('BWB001','6','TYR00006','SN00006','DOT-2022-006','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',20.0,16.0,15.5,16.5,5200.0,88.0,'mounted','Unit 1 - Pos 6'),
    ('BWB001','7','TYR00007','SN00007','DOT-2022-007','Radial','16.00R25','TRIANGLE','TB 516S',25.0,21.0,20.5,21.5,3000.0,92.0,'mounted','Unit 1 - Pos 7'),
    ('BWB001','8','TYR00008','SN00008','DOT-2022-008','Radial','16.00R25','Giti','GAO802',22.0,18.5,17.5,19.5,4500.0,90.0,'mounted','Unit 1 - Pos 8'),
    -- Unit 2 (BWB002): 4 mounted + 4 spare - positions 1-4 mounted, 5-8 spare
    ('BWB002','1','TYR00011','SN00011','DOT-2022-011','Radial','16.00R25','TECHKING','ET919+',27.0,23.0,22.0,24.0,2800.0,90.0,'mounted','Unit 2 - Pos 1'),
    ('BWB002','2','TYR00012','SN00012','DOT-2022-012','Radial','16.00R25','UNINEST','TIBERUN 811',29.0,25.0,24.5,25.5,1800.0,95.0,'mounted','Unit 2 - Pos 2'),
    ('BWB002','3','TYR00013','SN00013','DOT-2022-013','Radial','16.00R25','HILO','B01NL',21.0,17.0,16.0,18.0,5800.0,86.0,'mounted','Unit 2 - Pos 3'),
    ('BWB002','4','TYR00014','SN00014','DOT-2022-014','Radial','16.00R25','BRIDGESTONE','VUT',16.0,12.0,11.5,12.5,8200.0,82.0,'mounted','Unit 2 - Pos 4'),
    -- Unit 2 spares (no unit_id)
    ('BWB002','Spare','TYR00015','SN00015','DOT-2022-015','Radial','16.00R25','SAKURA','SRS-01',38.0,35.0,34.5,35.5,300.0,98.0,'spare','Unit 2 Spare - Pos 5'),
    ('BWB002','Spare','TYR00016','SN00016','DOT-2022-016','Radial','16.00R25','TECHKING','VUT',40.0,38.0,37.5,38.5,100.0,100.0,'spare','Unit 2 Spare - Pos 6'),
    ('BWB002','Spare','TYR00017','SN00017','DOT-2022-017','Radial','16.00R25','Giti','GAO802',36.0,33.0,32.5,33.5,600.0,96.0,'spare','Unit 2 Spare - Pos 7'),
    ('BWB002','Spare','TYR00018','SN00018','DOT-2022-018','Radial','16.00R25','ADVANCE','V-LUG',37.0,34.0,33.5,34.5,450.0,97.0,'spare','Unit 2 Spare - Pos 8'),
    -- Unit 3 (BWB003): all 14 spare
    ('BWB003','Spare','TYR00021','SN00021','DOT-2022-021','Radial','16.00R25','UNINEST','TIBERUN 811',42.0,40.0,39.5,40.5,0.0,100.0,'spare','Unit 3 Spare - 1'),
    ('BWB003','Spare','TYR00022','SN00022','DOT-2022-022','Radial','16.00R25','TECHKING','ET919',41.0,39.0,38.5,39.5,100.0,100.0,'spare','Unit 3 Spare - 2'),
    ('BWB003','Spare','TYR00023','SN00023','DOT-2022-023','Radial','16.00R25','HILO','B01NL',39.0,37.0,36.5,37.5,200.0,98.0,'spare','Unit 3 Spare - 3'),
    ('BWB003','Spare','TYR00024','SN00024','DOT-2022-024','Radial','16.00R25','BRIDGESTONE','VUT',38.0,36.0,35.5,36.5,350.0,97.0,'spare','Unit 3 Spare - 4'),
    ('BWB003','Spare','TYR00025','SN00025','DOT-2022-025','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',40.0,38.0,37.5,38.5,150.0,100.0,'spare','Unit 3 Spare - 5'),
    ('BWB003','Spare','TYR00026','SN00026','DOT-2022-026','Radial','16.00R25','SAKURA','SRS-01',37.0,35.0,34.5,35.5,400.0,96.0,'spare','Unit 3 Spare - 6'),
    ('BWB003','Spare','TYR00027','SN00027','DOT-2022-027','Radial','16.00R25','Giti','GAO802',39.0,37.0,36.5,37.5,250.0,98.0,'spare','Unit 3 Spare - 7'),
    ('BWB003','Spare','TYR00028','SN00028','DOT-2022-028','Radial','16.00R25','TRIANGLE','TB 516S',41.0,39.0,38.5,39.5,80.0,100.0,'spare','Unit 3 Spare - 8'),
    ('BWB003','Spare','TYR00029','SN00029','DOT-2022-029','Radial','16.00R25','ADVANCE','V-LUG',35.0,33.0,32.5,33.5,600.0,95.0,'spare','Unit 3 Spare - 9'),
    ('BWB003','Spare','TYR00030','SN00030','DOT-2022-030','Radial','16.00R25','TECHKING','ET919+',36.0,34.0,33.5,34.5,500.0,96.0,'spare','Unit 3 Spare - 10'),
    ('BWB003','Spare','TYR00031','SN00031','DOT-2022-031','Radial','16.00R25','UNINEST','TIBERUN 811',38.0,36.0,35.5,36.5,350.0,97.0,'spare','Unit 3 Spare - 11'),
    ('BWB003','Spare','TYR00032','SN00032','DOT-2022-032','Radial','16.00R25','HILO','B01NL',34.0,32.0,31.5,32.5,800.0,94.0,'spare','Unit 3 Spare - 12'),
    ('BWB003','Spare','TYR00033','SN00033','DOT-2022-033','Radial','16.00R25','BRIDGESTONE','VUT',33.0,31.0,30.5,31.5,900.0,93.0,'spare','Unit 3 Spare - 13'),
    ('BWB003','Spare','TYR00034','SN00034','DOT-2022-034','Radial','16.00R25','SAKURA','SRS-01',35.0,33.0,32.5,33.5,700.0,95.0,'spare','Unit 3 Spare - 14'),

    -- ========== BEI TYRES (35 total) ==========
    -- Unit 1 (BEI001): 8 mounted
    ('BEI001','1','TYR00051','SN00051','DOT-2023-001','Radial','16.00R25','TECHKING','ET919',29.0,25.0,24.5,25.5,1700.0,95.0,'mounted','Unit 1 - Pos 1'),
    ('BEI001','2','TYR00052','SN00052','DOT-2023-002','Radial','16.00R25','UNINEST','TIBERUN 811',27.0,23.0,22.0,24.0,2500.0,90.0,'mounted','Unit 1 - Pos 2'),
    ('BEI001','3','TYR00053','SN00053','DOT-2023-003','Radial','16.00R25','HILO','B01NL',31.0,27.5,27.0,28.0,1200.0,95.0,'mounted','Unit 1 - Pos 3'),
    ('BEI001','4','TYR00054','SN00054','DOT-2023-004','Radial','16.00R25','BRIDGESTONE','VUT',17.0,13.5,12.5,14.5,7200.0,84.0,'mounted','Unit 1 - Pos 4'),
    ('BEI001','5','TYR00055','SN00055','DOT-2023-005','Radial','16.00R25','ADVANCE','V-LUG',14.0,10.0,9.5,10.5,9800.0,78.0,'mounted','Unit 1 - Pos 5'),
    ('BEI001','6','TYR00056','SN00056','DOT-2023-006','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',21.0,17.0,16.5,17.5,4800.0,88.0,'mounted','Unit 1 - Pos 6'),
    ('BEI001','7','TYR00057','SN00057','DOT-2023-007','Radial','16.00R25','TRIANGLE','TB 516S',26.0,22.0,21.0,23.0,3100.0,92.0,'mounted','Unit 1 - Pos 7'),
    ('BEI001','8','TYR00058','SN00058','DOT-2023-008','Radial','16.00R25','Giti','GAO802',23.0,19.0,18.5,19.5,4100.0,90.0,'mounted','Unit 1 - Pos 8'),
    -- Unit 2 (BEI002): 4 mounted + 4 spare
    ('BEI002','1','TYR00061','SN00061','DOT-2023-011','Radial','16.00R25','TECHKING','ET919+',28.0,24.0,23.5,24.5,2400.0,90.0,'mounted','Unit 2 - Pos 1'),
    ('BEI002','2','TYR00062','SN00062','DOT-2023-012','Radial','16.00R25','UNINEST','TIBERUN 811',30.0,26.0,25.5,26.5,1500.0,95.0,'mounted','Unit 2 - Pos 2'),
    ('BEI002','3','TYR00063','SN00063','DOT-2023-013','Radial','16.00R25','HILO','B01NL',22.0,18.0,17.0,19.0,5400.0,86.0,'mounted','Unit 2 - Pos 3'),
    ('BEI002','4','TYR00064','SN00064','DOT-2023-014','Radial','16.00R25','BRIDGESTONE','VUT',15.0,11.0,10.5,11.5,8600.0,80.0,'mounted','Unit 2 - Pos 4'),
    ('BEI002','Spare','TYR00065','SN00065','DOT-2023-015','Radial','16.00R25','SAKURA','SRS-01',39.0,37.0,36.5,37.5,200.0,98.0,'spare','Unit 2 Spare - 5'),
    ('BEI002','Spare','TYR00066','SN00066','DOT-2023-016','Radial','16.00R25','TECHKING','VUT',41.0,39.0,38.5,39.5,50.0,100.0,'spare','Unit 2 Spare - 6'),
    ('BEI002','Spare','TYR00067','SN00067','DOT-2023-017','Radial','16.00R25','Giti','GAO802',38.0,36.0,35.5,36.5,400.0,97.0,'spare','Unit 2 Spare - 7'),
    ('BEI002','Spare','TYR00068','SN00068','DOT-2023-018','Radial','16.00R25','ADVANCE','V-LUG',37.0,35.0,34.5,35.5,550.0,96.0,'spare','Unit 2 Spare - 8'),
    -- Unit 3 (BEI003): 19 spare
    ('BEI003','Spare','TYR00071','SN00071','DOT-2023-021','Radial','16.00R25','UNINEST','TIBERUN 811',43.0,41.0,40.5,41.5,0.0,100.0,'spare','Unit 3 Spare - 1'),
    ('BEI003','Spare','TYR00072','SN00072','DOT-2023-022','Radial','16.00R25','TECHKING','ET919',42.0,40.0,39.5,40.5,80.0,100.0,'spare','Unit 3 Spare - 2'),
    ('BEI003','Spare','TYR00073','SN00073','DOT-2023-023','Radial','16.00R25','HILO','B01NL',40.0,38.0,37.5,38.5,150.0,98.0,'spare','Unit 3 Spare - 3'),
    ('BEI003','Spare','TYR00074','SN00074','DOT-2023-024','Radial','16.00R25','BRIDGESTONE','VUT',39.0,37.0,36.5,37.5,250.0,97.0,'spare','Unit 3 Spare - 4'),
    ('BEI003','Spare','TYR00075','SN00075','DOT-2023-025','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',41.0,39.0,38.5,39.5,100.0,100.0,'spare','Unit 3 Spare - 5'),
    ('BEI003','Spare','TYR00076','SN00076','DOT-2023-026','Radial','16.00R25','SAKURA','SRS-01',38.0,36.0,35.5,36.5,350.0,96.0,'spare','Unit 3 Spare - 6'),
    ('BEI003','Spare','TYR00077','SN00077','DOT-2023-027','Radial','16.00R25','Giti','GAO802',40.0,38.0,37.5,38.5,180.0,98.0,'spare','Unit 3 Spare - 7'),
    ('BEI003','Spare','TYR00078','SN00078','DOT-2023-028','Radial','16.00R25','TRIANGLE','TB 516S',42.0,40.0,39.5,40.5,60.0,100.0,'spare','Unit 3 Spare - 8'),
    ('BEI003','Spare','TYR00079','SN00079','DOT-2023-029','Radial','16.00R25','ADVANCE','V-LUG',36.0,34.0,33.5,34.5,500.0,95.0,'spare','Unit 3 Spare - 9'),
    ('BEI003','Spare','TYR00080','SN00080','DOT-2023-030','Radial','16.00R25','TECHKING','ET919+',37.0,35.0,34.5,35.5,400.0,96.0,'spare','Unit 3 Spare - 10'),
    ('BEI003','Spare','TYR00081','SN00081','DOT-2023-031','Radial','16.00R25','UNINEST','TIBERUN 811',39.0,37.0,36.5,37.5,300.0,97.0,'spare','Unit 3 Spare - 11'),
    ('BEI003','Spare','TYR00082','SN00082','DOT-2023-032','Radial','16.00R25','HILO','B01NL',35.0,33.0,32.5,33.5,700.0,94.0,'spare','Unit 3 Spare - 12'),
    ('BEI003','Spare','TYR00083','SN00083','DOT-2023-033','Radial','16.00R25','BRIDGESTONE','VUT',34.0,32.0,31.5,32.5,850.0,93.0,'spare','Unit 3 Spare - 13'),
    ('BEI003','Spare','TYR00084','SN00084','DOT-2023-034','Radial','16.00R25','SAKURA','SRS-01',36.0,34.0,33.5,34.5,600.0,95.0,'spare','Unit 3 Spare - 14'),
    ('BEI003','Spare','TYR00085','SN00085','DOT-2023-035','Radial','16.00R25','Giti','GAO802',38.0,36.0,35.5,36.5,450.0,97.0,'spare','Unit 3 Spare - 15'),
    ('BEI003','Spare','TYR00086','SN00086','DOT-2023-036','Radial','16.00R25','TECHKING','VUT',37.0,35.0,34.5,35.5,550.0,96.0,'spare','Unit 3 Spare - 16'),
    ('BEI003','Spare','TYR00087','SN00087','DOT-2023-037','Radial','16.00R25','ADVANCE','V-LUG',33.0,31.0,30.5,31.5,950.0,92.0,'spare','Unit 3 Spare - 17'),
    ('BEI003','Spare','TYR00088','SN00088','DOT-2023-038','Radial','16.00R25','TRIANGLE','TB 516S',35.0,33.0,32.5,33.5,750.0,95.0,'spare','Unit 3 Spare - 18'),
    ('BEI003','Spare','TYR00089','SN00089','DOT-2023-039','Radial','16.00R25','UNINEST','TIBERUN 811',40.0,38.0,37.5,38.5,200.0,98.0,'spare','Unit 3 Spare - 19'),

    -- ========== KPP TYRES (35 total) ==========
    -- Unit 1 (KPP001): 8 mounted
    ('KPP001','1','TYR00101','SN00101','DOT-2024-001','Radial','16.00R25','TECHKING','ET919',30.0,26.0,25.5,26.5,1300.0,95.0,'mounted','Unit 1 - Pos 1'),
    ('KPP001','2','TYR00102','SN00102','DOT-2024-002','Radial','16.00R25','UNINEST','TIBERUN 811',28.0,24.0,23.0,25.0,2100.0,90.0,'mounted','Unit 1 - Pos 2'),
    ('KPP001','3','TYR00103','SN00103','DOT-2024-003','Radial','16.00R25','HILO','B01NL',32.0,28.5,28.0,29.0,800.0,95.0,'mounted','Unit 1 - Pos 3'),
    ('KPP001','4','TYR00104','SN00104','DOT-2024-004','Radial','16.00R25','BRIDGESTONE','VUT',19.0,15.5,14.5,16.5,6200.0,85.0,'mounted','Unit 1 - Pos 4'),
    ('KPP001','5','TYR00105','SN00105','DOT-2024-005','Radial','16.00R25','ADVANCE','V-LUG',16.0,12.0,11.5,12.5,8700.0,80.0,'mounted','Unit 1 - Pos 5'),
    ('KPP001','6','TYR00106','SN00106','DOT-2024-006','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',22.0,18.0,17.5,18.5,4500.0,88.0,'mounted','Unit 1 - Pos 6'),
    ('KPP001','7','TYR00107','SN00107','DOT-2024-007','Radial','16.00R25','TRIANGLE','TB 516S',27.0,23.0,22.0,24.0,2800.0,92.0,'mounted','Unit 1 - Pos 7'),
    ('KPP001','8','TYR00108','SN00108','DOT-2024-008','Radial','16.00R25','Giti','GAO802',24.0,20.0,19.5,20.5,3800.0,90.0,'mounted','Unit 1 - Pos 8'),
    -- Unit 2 (KPP002): 4 mounted + 4 spare
    ('KPP002','1','TYR00111','SN00111','DOT-2024-011','Radial','16.00R25','TECHKING','ET919+',29.0,25.0,24.5,25.5,2000.0,90.0,'mounted','Unit 2 - Pos 1'),
    ('KPP002','2','TYR00112','SN00112','DOT-2024-012','Radial','16.00R25','UNINEST','TIBERUN 811',31.0,27.0,26.5,27.5,1200.0,95.0,'mounted','Unit 2 - Pos 2'),
    ('KPP002','3','TYR00113','SN00113','DOT-2024-013','Radial','16.00R25','HILO','B01NL',23.0,19.0,18.0,20.0,5100.0,86.0,'mounted','Unit 2 - Pos 3'),
    ('KPP002','4','TYR00114','SN00114','DOT-2024-014','Radial','16.00R25','BRIDGESTONE','VUT',17.0,13.0,12.5,13.5,7900.0,82.0,'mounted','Unit 2 - Pos 4'),
    ('KPP002','Spare','TYR00115','SN00115','DOT-2024-015','Radial','16.00R25','SAKURA','SRS-01',40.0,38.0,37.5,38.5,150.0,98.0,'spare','Unit 2 Spare - 5'),
    ('KPP002','Spare','TYR00116','SN00116','DOT-2024-016','Radial','16.00R25','TECHKING','VUT',42.0,40.0,39.5,40.5,30.0,100.0,'spare','Unit 2 Spare - 6'),
    ('KPP002','Spare','TYR00117','SN00117','DOT-2024-017','Radial','16.00R25','Giti','GAO802',39.0,37.0,36.5,37.5,350.0,97.0,'spare','Unit 2 Spare - 7'),
    ('KPP002','Spare','TYR00118','SN00118','DOT-2024-018','Radial','16.00R25','ADVANCE','V-LUG',38.0,36.0,35.5,36.5,500.0,96.0,'spare','Unit 2 Spare - 8'),
    -- Unit 3 (KPP003): 19 spare
    ('KPP003','Spare','TYR00121','SN00121','DOT-2024-021','Radial','16.00R25','UNINEST','TIBERUN 811',44.0,42.0,41.5,42.5,0.0,100.0,'spare','Unit 3 Spare - 1'),
    ('KPP003','Spare','TYR00122','SN00122','DOT-2024-022','Radial','16.00R25','TECHKING','ET919',43.0,41.0,40.5,41.5,50.0,100.0,'spare','Unit 3 Spare - 2'),
    ('KPP003','Spare','TYR00123','SN00123','DOT-2024-023','Radial','16.00R25','HILO','B01NL',41.0,39.0,38.5,39.5,120.0,98.0,'spare','Unit 3 Spare - 3'),
    ('KPP003','Spare','TYR00124','SN00124','DOT-2024-024','Radial','16.00R25','BRIDGESTONE','VUT',40.0,38.0,37.5,38.5,220.0,97.0,'spare','Unit 3 Spare - 4'),
    ('KPP003','Spare','TYR00125','SN00125','DOT-2024-025','Radial','16.00R25','VK TYRE','XTRA LOAD GRIP',42.0,40.0,39.5,40.5,80.0,100.0,'spare','Unit 3 Spare - 5'),
    ('KPP003','Spare','TYR00126','SN00126','DOT-2024-026','Radial','16.00R25','SAKURA','SRS-01',39.0,37.0,36.5,37.5,320.0,96.0,'spare','Unit 3 Spare - 6'),
    ('KPP003','Spare','TYR00127','SN00127','DOT-2024-027','Radial','16.00R25','Giti','GAO802',41.0,39.0,38.5,39.5,150.0,98.0,'spare','Unit 3 Spare - 7'),
    ('KPP003','Spare','TYR00128','SN00128','DOT-2024-028','Radial','16.00R25','TRIANGLE','TB 516S',43.0,41.0,40.5,41.5,40.0,100.0,'spare','Unit 3 Spare - 8'),
    ('KPP003','Spare','TYR00129','SN00129','DOT-2024-029','Radial','16.00R25','ADVANCE','V-LUG',37.0,35.0,34.5,35.5,480.0,95.0,'spare','Unit 3 Spare - 9'),
    ('KPP003','Spare','TYR00130','SN00130','DOT-2024-030','Radial','16.00R25','TECHKING','ET919+',38.0,36.0,35.5,36.5,380.0,96.0,'spare','Unit 3 Spare - 10'),
    ('KPP003','Spare','TYR00131','SN00131','DOT-2024-031','Radial','16.00R25','UNINEST','TIBERUN 811',40.0,38.0,37.5,38.5,280.0,97.0,'spare','Unit 3 Spare - 11'),
    ('KPP003','Spare','TYR00132','SN00132','DOT-2024-032','Radial','16.00R25','HILO','B01NL',36.0,34.0,33.5,34.5,680.0,94.0,'spare','Unit 3 Spare - 12'),
    ('KPP003','Spare','TYR00133','SN00133','DOT-2024-033','Radial','16.00R25','BRIDGESTONE','VUT',35.0,33.0,32.5,33.5,820.0,93.0,'spare','Unit 3 Spare - 13'),
    ('KPP003','Spare','TYR00134','SN00134','DOT-2024-034','Radial','16.00R25','SAKURA','SRS-01',37.0,35.0,34.5,35.5,580.0,95.0,'spare','Unit 3 Spare - 14'),
    ('KPP003','Spare','TYR00135','SN00135','DOT-2024-035','Radial','16.00R25','Giti','GAO802',39.0,37.0,36.5,37.5,420.0,97.0,'spare','Unit 3 Spare - 15'),
    ('KPP003','Spare','TYR00136','SN00136','DOT-2024-036','Radial','16.00R25','TECHKING','VUT',38.0,36.0,35.5,36.5,530.0,96.0,'spare','Unit 3 Spare - 16'),
    ('KPP003','Spare','TYR00137','SN00137','DOT-2024-037','Radial','16.00R25','ADVANCE','V-LUG',34.0,32.0,31.5,32.5,920.0,92.0,'spare','Unit 3 Spare - 17'),
    ('KPP003','Spare','TYR00138','SN00138','DOT-2024-038','Radial','16.00R25','TRIANGLE','TB 516S',36.0,34.0,33.5,34.5,720.0,95.0,'spare','Unit 3 Spare - 18'),
    ('KPP003','Spare','TYR00139','SN00139','DOT-2024-039','Radial','16.00R25','UNINEST','TIBERUN 811',41.0,39.0,38.5,39.5,180.0,98.0,'spare','Unit 3 Spare - 19')
) AS t(unit_code,mounted_pos,barcode,serial,dot,tyre_type,size_name,brand_name,pattern_name,otd,rtd,rtd1,rtd2,lifetime,psi,status,remarks)
LEFT JOIN master_sizes sz ON sz.name = t.size_name
LEFT JOIN master_brands b ON b.name = t.brand_name
LEFT JOIN master_patterns pt ON pt.name = t.pattern_name AND pt.brand_id = b.id
WHERE (c.name = 'Berkat Anuegrah Sejahtera' AND t.unit_code IN ('BWB001','BWB002','BWB003'))
   OR (c.name = 'Borneo Energy Indonesia' AND t.unit_code IN ('BEI001','BEI002','BEI003'))
   OR (c.name = 'Kalimantan Prima Persada' AND t.unit_code IN ('KPP001','KPP002','KPP003'))
ON CONFLICT DO NOTHING;
