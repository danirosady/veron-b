-- Rollback: remove seeded company data (order matters — FK constraints)
DELETE FROM replacement_details WHERE replacement_id IN (
    SELECT r.id FROM replacements r
    JOIN companies c ON r.company_id = c.id WHERE c.name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM replacements WHERE company_id IN (
    SELECT id FROM companies WHERE name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM tyre_master WHERE company_id IN (
    SELECT id FROM companies WHERE name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM units WHERE company_id IN (
    SELECT id FROM companies WHERE name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM drivers WHERE company_id IN (
    SELECT id FROM companies WHERE name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM projects WHERE company_id IN (
    SELECT id FROM companies WHERE name = 'Berkat Anuegrah Sejahtera'
);
DELETE FROM companies WHERE name = 'Berkat Anuegrah Sejahtera';
