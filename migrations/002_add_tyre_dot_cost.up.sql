-- ============================================================
-- Migration: 002_add_tyre_dot_cost
-- Add DOT code (short format) and Cost columns to tyre_master
-- ============================================================

-- Add DOT column (e.g., "3522" = week 35, 2022)
ALTER TABLE tyre_master ADD COLUMN dot VARCHAR(10);

-- Add Cost column (tyre price)
ALTER TABLE tyre_master ADD COLUMN cost DECIMAL(15,2);

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_tyre_master_dot ON tyre_master(dot);
CREATE INDEX IF NOT EXISTS idx_tyre_master_cost ON tyre_master(cost);
