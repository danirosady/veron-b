-- ============================================================
-- Rollback: 004_add_new_tyre_hm
-- ============================================================

DROP INDEX IF EXISTS idx_replacement_details_new_tyre_hm;
ALTER TABLE replacement_details DROP COLUMN IF EXISTS new_tyre_hm;
