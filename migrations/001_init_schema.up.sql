-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- COMPANIES (created first — referenced by users, projects, drivers, units, tyres)
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
-- USERS (references companies)
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL DEFAULT 'admin_company',
    -- role: superadmin | admin_company
    company_id      BIGINT NULL,
    -- untuk admin_company, NULL untuk superadmin
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    -- status: active | inactive
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
-- PROJECTS (references companies)
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
-- DRIVERS (references companies)
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
-- MASTER BRANDS (no dependencies)
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
-- MASTER SIZES (no dependencies)
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
-- MASTER TYPES (no dependencies)
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
-- MASTER PATTERNS (references master_brands)
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
-- MASTER REASONS (no dependencies)
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
-- MASTER ACTIONS (no dependencies)
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
-- MASTER REMARKS (no dependencies)
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
-- UNIT TYPE CONFIGS (no dependencies)
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
-- UNITS (references companies, projects)
-- ============================================================
CREATE TABLE IF NOT EXISTS units (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    project_id         BIGINT NOT NULL,
    unit_id            VARCHAR(50) NOT NULL UNIQUE,
    -- kode unik kendaraan, contoh: BWB051, GD33
    unit_model         VARCHAR(255) NOT NULL,
    -- contoh: SANY SKT 105S, SANY 105S, GREADER
    plate_number       VARCHAR(50) NULL,
    tyre_size_default  VARCHAR(50) NOT NULL,
    -- contoh: 16.00R25, 14.00R24
    unit_type          VARCHAR(50) NOT NULL DEFAULT 'ADT',
    -- unit_type: SANY_10POS | GREADER_6POS | ADT_8POS
    max_position       INTEGER NOT NULL DEFAULT 6,
    -- jumlah posisi ban maksimum
    current_hm         DECIMAL(15,2) NOT NULL DEFAULT 0,
    -- Hour Meter saat ini
    status             VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_units_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT fk_units_project FOREIGN KEY (project_id)
        REFERENCES projects(id) ON DELETE RESTRICT,

    CONSTRAINT chk_units_max_position CHECK (max_position > 0 AND max_position <= 20),
    CONSTRAINT chk_units_status CHECK (status IN ('active', 'inactive'))
);

CREATE INDEX IF NOT EXISTS idx_units_company_id ON units(company_id);
CREATE INDEX IF NOT EXISTS idx_units_project_id ON units(project_id);
CREATE INDEX IF NOT EXISTS idx_units_unit_id ON units(unit_id);
CREATE INDEX IF NOT EXISTS idx_units_unit_type ON units(unit_type);
CREATE INDEX IF NOT EXISTS idx_units_status ON units(status);

-- ============================================================
-- TYRE MASTER (references companies, units, master tables)
-- ============================================================
CREATE TABLE IF NOT EXISTS tyre_master (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    unit_id             BIGINT NULL,
    -- NULL = spare / tidak terpasang
    mounted_position    INTEGER NULL,
    -- posisi ban di unit, NULL jika tidak terpasang

    -- Identifikasi
    barcode             VARCHAR(100) NOT NULL UNIQUE,
    -- barcode scanner
    serial_number       VARCHAR(100) NOT NULL UNIQUE,
    -- S/N TYRE
    dot_code            VARCHAR(100) NULL,
    -- Department of Transportation

    -- Spesifikasi
    type                VARCHAR(50) NOT NULL DEFAULT 'Radial',
    size_id             BIGINT NOT NULL,
    brand_id           BIGINT NOT NULL,
    pattern_id          BIGINT NOT NULL,

    -- Kondisi ban
    otd                 DECIMAL(5,2) NOT NULL DEFAULT 0,
    rtd                 DECIMAL(5,2) NOT NULL DEFAULT 0,
    rtd_1               DECIMAL(5,2) NULL,
    rtd_2               DECIMAL(5,2) NULL,
    lifetime             DECIMAL(15,2) NOT NULL DEFAULT 0,
    psi                  DECIMAL(5,2) NULL,

    -- Status
    status              VARCHAR(20) NOT NULL DEFAULT 'spare',
    -- status: spare | mounted | dismounted | scrap
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
-- REPLACEMENTS (references companies, projects, units, drivers, users)
-- ============================================================
CREATE TABLE IF NOT EXISTS replacements (
    id                  BIGSERIAL PRIMARY KEY,
    company_id          BIGINT NOT NULL,
    project_id          BIGINT NOT NULL,
    unit_id             BIGINT NOT NULL,
    driver_id           BIGINT NOT NULL,
    date                DATE NOT NULL,
    hm_update           DECIMAL(15,2) NOT NULL DEFAULT 0,
    -- HM saat penggantian terakhir
    current_life_hm     DECIMAL(15,2) NOT NULL,
    -- HM unit saat ini
    hm_plan             DECIMAL(15,2) NOT NULL,
    -- rencana HM penggantian berikutnya
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
-- REPLACEMENT DETAILS (references replacements, tyre_master, master_reasons)
-- ============================================================
CREATE TABLE IF NOT EXISTS replacement_details (
    id                          BIGSERIAL PRIMARY KEY,
    replacement_id              BIGINT NOT NULL,

    -- Posisi & Aksi
    position                    INTEGER NOT NULL,
    action                      VARCHAR(50) NOT NULL,
    -- action: mount | dismount | swap

    -- Old Tyre
    old_tyre_id                BIGINT NULL,
    old_tyre_serial_number     VARCHAR(100) NULL,
    old_tyre_pattern           VARCHAR(100) NULL,
    old_tyre_size              VARCHAR(50) NULL,
    old_tyre_tread_1          DECIMAL(5,2) NULL,
    old_tyre_tread_2          DECIMAL(5,2) NULL,
    old_tyre_lifetime          DECIMAL(15,2) NULL,
    old_tyre_status            VARCHAR(20) NULL,
    failure_reason_id          BIGINT NULL,
    from_unit_id               VARCHAR(50) NULL,

    -- New Tyre
    new_tyre_id                BIGINT NULL,
    new_tyre_serial_number     VARCHAR(100) NULL,
    new_tyre_pattern           VARCHAR(100) NULL,
    new_tyre_size             VARCHAR(50) NULL,
    new_tyre_tread_1          DECIMAL(5,2) NULL,
    new_tyre_tread_2          DECIMAL(5,2) NULL,
    new_tyre_current_lifetime DECIMAL(15,2) NOT NULL DEFAULT 0,
    new_tyre_status           VARCHAR(50) NULL,

    -- Output
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
