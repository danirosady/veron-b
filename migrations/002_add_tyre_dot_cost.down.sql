-- ============================================================
-- Migration: 002_add_tyre_dot_cost (down)
-- Rollback: Remove DOT and Cost columns from tyre_master
-- ============================================================

-- Drop indexes
DROP INDEX IF EXISTS idx_tyre_master_dot;
DROP INDEX IF EXISTS idx_tyre_master_cost;

-- Remove columns
ALTER TABLE tyre_master DROP COLUMN IF EXISTS dot;
ALTER TABLE tyre_master DROP COLUMN IF EXISTS cost;
