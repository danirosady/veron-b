-- ============================================================
-- Migration: 005_add_employees (down)
-- Restores drivers table and removes employees table
-- ============================================================

-- Drop junction table
DROP TABLE IF EXISTS replacement_employees CASCADE;

-- Drop supervisor_id column
ALTER TABLE replacement_details DROP CONSTRAINT IF EXISTS fk_rd_supervisor;
ALTER TABLE replacement_details DROP COLUMN IF EXISTS supervisor_id;

-- Restore drivers table
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

-- Copy driver role employees back to drivers table
INSERT INTO drivers (company_id, name, employee_id, phone, status, created_at, updated_at)
SELECT company_id, name, employee_id, phone, status, created_at, updated_at
FROM employees WHERE role = 'driver';

-- Restore replacement.driver_id FK to drivers
ALTER TABLE replacements DROP CONSTRAINT IF EXISTS fk_replacements_driver;
ALTER TABLE replacements ADD COLUMN driver_driver_id BIGINT;
UPDATE replacements SET driver_driver_id = driver_id;
ALTER TABLE replacements DROP COLUMN driver_id;
ALTER TABLE replacements RENAME COLUMN driver_driver_id TO driver_id;
ALTER TABLE replacements ADD CONSTRAINT fk_replacements_driver
    FOREIGN KEY (driver_id) REFERENCES drivers(id) ON DELETE RESTRICT;

-- Drop employees table
DROP TABLE IF EXISTS employees CASCADE;
