-- ============================================================
-- Rollback: 006_rename_master_remarks_and_add_tyre_remarks_fk
-- ============================================================

-- 1. Remove FK index
DROP INDEX IF EXISTS idx_tyre_master_remarks_id;

-- 2. Remove remarks_id column (this also removes the FK constraint)
ALTER TABLE tyre_master DROP COLUMN IF EXISTS remarks_id;

-- 3. Re-add free-text remarks column
ALTER TABLE tyre_master ADD COLUMN remarks TEXT;

-- 4. Restore CHECK constraint to original (without 'new_tyre')
ALTER TABLE tyre_master DROP CONSTRAINT IF EXISTS chk_tyre_master_status;
ALTER TABLE tyre_master ADD CONSTRAINT chk_tyre_master_status
    CHECK (status IN ('spare', 'mounted', 'dismounted', 'scrap'));

-- 5. Restore master_remarks.name column
ALTER TABLE master_remarks RENAME COLUMN status_name TO name;

-- 6. Clean up extra status_name rows inserted by up migration
DELETE FROM master_remarks WHERE name IN (
    'Mounted', 'Dismounted'
);
