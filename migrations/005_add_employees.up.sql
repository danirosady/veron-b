-- ============================================================
-- Migration: 005_add_employees
-- Replaces drivers table with unified employees table
-- Adds supervisor_id to replacement_details
-- Adds replacement_employees junction for tyreman
-- ============================================================

-- Create employees table (replaces drivers)
CREATE TABLE IF NOT EXISTS employees (
    id              BIGSERIAL PRIMARY KEY,
    company_id      BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    employee_id     VARCHAR(50) NOT NULL,
    phone           VARCHAR(50) NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'driver',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_employees_company FOREIGN KEY (company_id)
        REFERENCES companies(id) ON DELETE RESTRICT,
    CONSTRAINT uq_employees_company_employee_role UNIQUE (company_id, employee_id, role),
    CONSTRAINT chk_employees_role CHECK (role IN ('driver', 'supervisor', 'tyreman')),
    CONSTRAINT chk_employees_status CHECK (status IN ('active', 'inactive'))
);
CREATE INDEX IF NOT EXISTS idx_employees_company_id ON employees(company_id);
CREATE INDEX IF NOT EXISTS idx_employees_role ON employees(role);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);

-- Copy existing driver data to employees table
INSERT INTO employees (company_id, name, employee_id, phone, role, status, created_at, updated_at)
SELECT company_id, name, employee_id, phone, 'driver', status, created_at, updated_at
FROM drivers;

-- Seed Supervisors and Tyremen (2 each per company)
INSERT INTO employees (company_id, name, employee_id, phone, role, status)
SELECT c.id, e.name, e.emp_id, e.phone, e.role, 'active'
FROM companies c
CROSS JOIN (VALUES
    -- BAS
    ('Budi Santoso', 'SPV001', '0813-1111-3333', 'supervisor'),
    ('Dewi Lestari', 'SPV002', '0813-2222-4444', 'supervisor'),
    ('Junaedi', 'TYM001', '0813-3333-5555', 'tyreman'),
    ('Kartika Sari', 'TYM002', '0813-4444-6666', 'tyreman'),
    -- BEI
    ('Ahmad Fauzi', 'SPV010', '0813-5555-7777', 'supervisor'),
    ('Nurul Huda', 'SPV011', '0813-6666-8888', 'supervisor'),
    ('Siti Aminah', 'TYM010', '0813-7777-9999', 'tyreman'),
    ('Wawan Setiawan', 'TYM011', '0813-8888-0000', 'tyreman'),
    -- KPP
    ('Fathan Rizki', 'SPV020', '0813-0000-1111', 'supervisor'),
    ('Mega Putri', 'SPV021', '0813-1111-2222', 'supervisor'),
    ('Dimas Arya', 'TYM020', '0813-2222-3333', 'tyreman'),
    ('Lina Marlina', 'TYM021', '0813-3333-4444', 'tyreman')
) AS e(name, emp_id, phone, role)
WHERE (c.name = 'Berkat Anuegrah Sejahtera' AND e.emp_id IN ('SPV001','SPV002','TYM001','TYM002'))
   OR (c.name = 'Borneo Energy Indonesia' AND e.emp_id IN ('SPV010','SPV011','TYM010','TYM011'))
   OR (c.name = 'Kalimantan Prima Persada' AND e.emp_id IN ('SPV020','SPV021','TYM020','TYM021'))
ON CONFLICT (company_id, employee_id, role) DO NOTHING;

-- Migrate replacements.driver_id -> employees.id
-- First add temporary column
ALTER TABLE replacements ADD COLUMN driver_employee_id BIGINT;
-- Map old driver IDs to new employee IDs (same ID but now in employees table)
UPDATE replacements SET driver_employee_id = driver_id;
-- Drop old FK and column
ALTER TABLE replacements DROP CONSTRAINT IF EXISTS fk_replacements_driver;
ALTER TABLE replacements DROP COLUMN driver_id;
-- Rename temp column
ALTER TABLE replacements RENAME COLUMN driver_employee_id TO driver_id;
-- Add FK back pointing to employees
ALTER TABLE replacements ADD CONSTRAINT fk_replacements_driver
    FOREIGN KEY (driver_id) REFERENCES employees(id) ON DELETE RESTRICT;

-- Add supervisor_id to replacement_details
ALTER TABLE replacement_details ADD COLUMN supervisor_id BIGINT;
ALTER TABLE replacement_details ADD CONSTRAINT fk_rd_supervisor
    FOREIGN KEY (supervisor_id) REFERENCES employees(id) ON DELETE SET NULL;

-- Create junction table for tyreman (many-to-many: replacement_detail can have multiple tyremen)
CREATE TABLE IF NOT EXISTS replacement_employees (
    id                      BIGSERIAL PRIMARY KEY,
    replacement_detail_id   BIGINT NOT NULL,
    employee_id             BIGINT NOT NULL,
    created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_re_Replacement_detail FOREIGN KEY (replacement_detail_id)
        REFERENCES replacement_details(id) ON DELETE CASCADE,
    CONSTRAINT fk_re_employee FOREIGN KEY (employee_id)
        REFERENCES employees(id) ON DELETE RESTRICT,
    CONSTRAINT uq_re_replacement_detail_employee UNIQUE (replacement_detail_id, employee_id)
);
CREATE INDEX IF NOT EXISTS idx_replacement_employees_detail_id ON replacement_employees(replacement_detail_id);
CREATE INDEX IF NOT EXISTS idx_replacement_employees_employee_id ON replacement_employees(employee_id);

-- Drop old drivers table
DROP TABLE IF EXISTS drivers CASCADE;
