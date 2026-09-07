-- Rollback: restore old position_config to original seed format
-- axle: poros_1→front, poros_2→bogie, poros_3→rear_1, poros_4→rear_2, poros_5→rear
-- side: left/right → front_left/front_right, rear_left/rear_right etc.
-- position: string number → "Pos N"
-- label: "Tyre N" → original descriptive label

-- ── ADT_8POS ────────────────────────────────────────────────────────────────
UPDATE unit_type_configs
SET position_config = '[{"position": "Pos 1", "label": "R-L1", "side": "rear_left",   "axle": "rear_1", "x": 0.15, "y": 0.25, "mirror_of": "Pos 2"},
  {"position": "Pos 2", "label": "R-R1", "side": "rear_right",  "axle": "rear_1", "x": 0.85, "y": 0.25, "mirror_of": "Pos 1"},
  {"position": "Pos 3", "label": "R-L2", "side": "rear_left",   "axle": "rear_2", "x": 0.25, "y": 0.25, "mirror_of": "Pos 4"},
  {"position": "Pos 4", "label": "R-R2", "side": "rear_right",  "axle": "rear_2", "x": 0.75, "y": 0.25, "mirror_of": "Pos 3"},
  {"position": "Pos 5", "label": "R-L3", "side": "rear_left",   "axle": "bogie",   "x": 0.20, "y": 0.55, "mirror_of": "Pos 6"},
  {"position": "Pos 6", "label": "R-R3", "side": "rear_right",  "axle": "bogie",   "x": 0.80, "y": 0.55, "mirror_of": "Pos 5"},
  {"position": "Pos 7", "label": "F-L",  "side": "front_left",  "axle": "front",   "x": 0.35, "y": 0.75, "mirror_of": "Pos 8"},
  {"position": "Pos 8", "label": "F-R",  "side": "front_right", "axle": "front",   "x": 0.65, "y": 0.75, "mirror_of": "Pos 7"}]'::jsonb
WHERE unit_type = 'ADT_8POS';

-- ── SANY_10POS ─────────────────────────────────────────────────────────────
UPDATE unit_type_configs
SET position_config = '[{"position": "Pos 1", "label": "R-L1", "side": "rear_left",   "axle": "rear_1", "x": 0.10, "y": 0.25, "mirror_of": "Pos 2"},
  {"position": "Pos 2", "label": "R-L2", "side": "rear_left",   "axle": "rear_1", "x": 0.25, "y": 0.25, "mirror_of": "Pos 1"},
  {"position": "Pos 3", "label": "R-L3", "side": "rear_left",   "axle": "rear_2", "x": 0.40, "y": 0.25, "mirror_of": "Pos 4"},
  {"position": "Pos 4", "label": "R-R1", "side": "rear_right",  "axle": "rear_1", "x": 0.60, "y": 0.25, "mirror_of": "Pos 3"},
  {"position": "Pos 5", "label": "R-R2", "side": "rear_right",  "axle": "rear_2", "x": 0.75, "y": 0.25, "mirror_of": "Pos 3"},
  {"position": "Pos 6", "label": "R-R3", "side": "rear_right",  "axle": "rear_2", "x": 0.90, "y": 0.25, "mirror_of": "Pos 5"},
  {"position": "Pos 7", "label": "R-L4", "side": "rear_left",   "axle": "bogie",   "x": 0.17, "y": 0.55, "mirror_of": "Pos 8"},
  {"position": "Pos 8", "label": "R-R4", "side": "rear_right",  "axle": "bogie",   "x": 0.83, "y": 0.55, "mirror_of": "Pos 7"},
  {"position": "Pos 9", "label": "F-L1", "side": "front_left",  "axle": "front",   "x": 0.35, "y": 0.75, "mirror_of": "Pos 10"},
  {"position": "Pos 10","label": "F-R1", "side": "front_right", "axle": "front",   "x": 0.65, "y": 0.75, "mirror_of": "Pos 9"}]'::jsonb
WHERE unit_type = 'SANY_10POS';

-- ── GREADER_6POS ───────────────────────────────────────────────────────────
UPDATE unit_type_configs
SET position_config = '[{"position": "Pos 1", "label": "R-L",  "side": "rear_left",   "axle": "rear",   "x": 0.20, "y": 0.30, "mirror_of": "Pos 2"},
  {"position": "Pos 2", "label": "R-R",  "side": "rear_right",  "axle": "rear",   "x": 0.80, "y": 0.30, "mirror_of": "Pos 1"},
  {"position": "Pos 3", "label": "F-L",  "side": "front_left",  "axle": "front",   "x": 0.30, "y": 0.70, "mirror_of": "Pos 4"},
  {"position": "Pos 4", "label": "F-R",  "side": "front_right", "axle": "front",   "x": 0.70, "y": 0.70, "mirror_of": "Pos 3"},
  {"position": "Pos 5", "label": "R-L2", "side": "rear_left",   "axle": "rear_2", "x": 0.20, "y": 0.55, "mirror_of": "Pos 6"},
  {"position": "Pos 6", "label": "R-R2", "side": "rear_right",  "axle": "rear_2", "x": 0.80, "y": 0.55, "mirror_of": "Pos 5"}]'::jsonb
WHERE unit_type = 'GREADER_6POS';
