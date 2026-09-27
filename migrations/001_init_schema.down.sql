-- ============================================================
-- Rollback: 001_init_schema
-- Truncate all tables in reverse dependency order (child → parent)
-- RESTART IDENTITY to reset serial counters
-- ============================================================

TRUNCATE
    replacement_details,
    replacements,
    tyre_master,
    units,
    drivers,
    projects,
    companies
RESTART IDENTITY CASCADE;

TRUNCATE
    unit_type_configs,
    master_patterns,
    master_reasons,
    master_actions,
    master_remarks,
    master_types,
    master_sizes,
    master_brands,
    users
RESTART IDENTITY CASCADE;
