-- ============================================================
-- Migration: 006_rename_master_remarks_and_add_tyre_remarks_fk
-- Changes:
--   1. Rename master_remarks.name -> master_remarks.status_name
--   2. Add remarks_id FK column to tyre_master (references master_remarks.id)
--   3. Drop old free-text remarks column from tyre_master
--   4. Update CHECK constraint to include 'new_tyre' status
--   5. Update master_remarks seed with complete status list + status_name column
--   6. Backfill remarks_id for existing tyres based on current status
-- ============================================================

-- 1. Rename master_remarks.name -> status_name
ALTER TABLE master_remarks RENAME COLUMN name TO status_name;

-- 2. Add remarks_id column to tyre_master
ALTER TABLE tyre_master ADD COLUMN remarks_id BIGINT REFERENCES master_remarks(id) ON DELETE SET NULL;

-- 3. Drop old free-text remarks column (no longer needed - replaced by FK)
ALTER TABLE tyre_master DROP COLUMN remarks;

-- 4. Update CHECK constraint on tyre_master.status to include 'new_tyre'
ALTER TABLE tyre_master DROP CONSTRAINT IF EXISTS chk_tyre_master_status;
ALTER TABLE tyre_master ADD CONSTRAINT chk_tyre_master_status
    CHECK (status IN ('new_tyre', 'spare', 'mounted', 'dismounted', 'scrap'));

-- 5. Update master_remarks seed - replace old names with status_name column
--    First delete old seed entries, then insert fresh ones with explicit IDs
DELETE FROM master_remarks WHERE status_name IN (
    'New Tyre', 'Repair', 'Running', 'Spare', 'Scrap', 'UNIT SECOND HAND', 'Ex Repair', 'Canibal Tyre'
);

-- Insert all tyre status values as master_remarks with FIXED IDs
-- These IDs are used in step 6 for backfill, so they must be stable
INSERT INTO master_remarks (id, status_name, status) VALUES
    (1, 'New Tyre',       'active'),
    (2, 'Repair',         'active'),
    (3, 'Running',        'active'),
    (4, 'Spare',          'active'),
    (5, 'Scrap',          'active'),
    (6, 'UNIT SECOND HAND','active'),
    (7, 'Ex Repair',      'active'),
    (8, 'Canibal Tyre',   'active'),
    (9, 'Mounted',        'active'),
    (10,'Dismounted',     'active')
ON CONFLICT DO NOTHING;

-- Reset sequence so next insert uses correct ID after gaps
SELECT setval('master_remarks_id_seq', (SELECT COALESCE(MAX(id), 0) FROM master_remarks));

-- 6. Backfill remarks_id for existing tyres based on current status
--    This maps tyre_master.status -> master_remarks.status_name
UPDATE tyre_master
SET remarks_id = mr.id
FROM master_remarks mr
WHERE mr.status_name = CASE tyre_master.status
    WHEN 'mounted'    THEN 'Mounted'
    WHEN 'dismounted' THEN 'Dismounted'
    WHEN 'spare'     THEN 'Spare'
    WHEN 'scrap'     THEN 'Scrap'
    WHEN 'new_tyre'  THEN 'New Tyre'
    ELSE NULL
END;

-- Create index for fast FK lookups
CREATE INDEX IF NOT EXISTS idx_tyre_master_remarks_id ON tyre_master(remarks_id);
