-- ============================================================
-- Migration: 004_add_new_tyre_hm
-- Add new_tyre_hm column to replacement_details for HM tracking at mount
-- Used to calculate old_tyre_lifetime = current_life_hm - new_tyre_hm on dismount
-- ============================================================

ALTER TABLE replacement_details ADD COLUMN IF NOT EXISTS new_tyre_hm DECIMAL(15,2);

CREATE INDEX IF NOT EXISTS idx_replacement_details_new_tyre_hm ON replacement_details(new_tyre_hm);
