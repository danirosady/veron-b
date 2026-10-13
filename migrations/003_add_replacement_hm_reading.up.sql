-- ============================================================
-- Migration: 003_add_replacement_hm_reading
-- Add hm_reading column to replacement_details for queue process
-- ============================================================

-- Add HM Reading column to capture HM at time of replacement
ALTER TABLE replacement_details ADD COLUMN IF NOT EXISTS hm_reading DECIMAL(15,2);

-- Create index for new column
CREATE INDEX IF NOT EXISTS idx_replacement_details_hm_reading ON replacement_details(hm_reading);
