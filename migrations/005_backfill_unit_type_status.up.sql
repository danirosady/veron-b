-- Backfill status column for unit_type_configs
-- Existing rows may have NULL status because the Go entity never had the Status field
UPDATE unit_type_configs SET status = 'active' WHERE status IS NULL;

-- Also ensure seed data includes status going forward (re-run-safe)
