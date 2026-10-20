-- ============================================================
-- Migration: 007_add_tyre_status_logs
-- Creates the audit log table for tyre status and remarks changes.
-- ============================================================

CREATE TABLE IF NOT EXISTS tyre_status_logs (
    id              BIGSERIAL PRIMARY KEY,
    tyre_id         BIGINT NOT NULL,
    old_status      VARCHAR(20) NULL,
    old_remarks_id  BIGINT NULL,
    new_status      VARCHAR(20) NOT NULL,
    new_remarks_id  BIGINT NULL,
    changed_by      BIGINT NULL,
    reason          TEXT NULL,
    changed_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tsl_tyre FOREIGN KEY (tyre_id)
        REFERENCES tyre_master(id) ON DELETE CASCADE,
    CONSTRAINT fk_tsl_old_remarks FOREIGN KEY (old_remarks_id)
        REFERENCES master_remarks(id) ON DELETE SET NULL,
    CONSTRAINT fk_tsl_new_remarks FOREIGN KEY (new_remarks_id)
        REFERENCES master_remarks(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_tyre_status_logs_tyre_id ON tyre_status_logs(tyre_id);
CREATE INDEX IF NOT EXISTS idx_tyre_status_logs_changed_at ON tyre_status_logs(changed_at);
