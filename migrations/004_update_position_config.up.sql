-- Migration: 004_update_position_config
-- Replaces old axle keys (front/bogie/rear_1/rear_2/rear) with poros_1-5,
-- old side values (front_left etc.) with bare left/right,
-- and regenerates labels from "Pos N" / descriptive to "Tyre N".
-- ADT_8POS: 8 positions across 4 axles (poros_3/poros_4/poros_2/poros_1)
-- SANY_10POS: 12 positions across 3 axles
-- GREADER_6POS: 6 positions across 3 axles

-- ── ADT_8POS (8 positions) ──────────────────────────────────────────────────
UPDATE unit_type_configs
SET position_config = '[
  {"position":"1","label":"Tyre 1","side":"left","axle":"poros_3","x":0.15,"y":0.25},
  {"position":"2","label":"Tyre 2","side":"right","axle":"poros_3","x":0.85,"y":0.25},
  {"position":"3","label":"Tyre 3","side":"left","axle":"poros_4","x":0.25,"y":0.25},
  {"position":"4","label":"Tyre 4","side":"right","axle":"poros_4","x":0.75,"y":0.25},
  {"position":"5","label":"Tyre 5","side":"left","axle":"poros_2","x":0.20,"y":0.55},
  {"position":"6","label":"Tyre 6","side":"right","axle":"poros_2","x":0.80,"y":0.55},
  {"position":"7","label":"Tyre 7","side":"left","axle":"poros_1","x":0.35,"y":0.75},
  {"position":"8","label":"Tyre 8","side":"right","axle":"poros_1","x":0.65,"y":0.75}
]'::jsonb,
    max_position = 8
WHERE unit_type = 'ADT_8POS';

-- ── SANY_10POS ─────────────────────────────────────────────────────────────
-- Old: rear_1/rear_2/bogie/front at y=0.25/0.25/0.55/0.75
-- New: poros_3/poros_4/poros_2/poros_1
UPDATE unit_type_configs
SET position_config = '[
  {"position": "1", "label": "Tyre 1", "side": "left",  "axle": "poros_3", "x": 0.10, "y": 0.25, "mirror_of": "2"},
  {"position": "2", "label": "Tyre 2", "side": "right", "axle": "poros_3", "x": 0.25, "y": 0.25, "mirror_of": "1"},
  {"position": "3", "label": "Tyre 3", "side": "left",  "axle": "poros_3", "x": 0.40, "y": 0.25, "mirror_of": "4"},
  {"position": "4", "label": "Tyre 4", "side": "right", "axle": "poros_3", "x": 0.60, "y": 0.25, "mirror_of": "3"},
  {"position": "5", "label": "Tyre 5", "side": "left",  "axle": "poros_4", "x": 0.40, "y": 0.25, "mirror_of": "6"},
  {"position": "6", "label": "Tyre 6", "side": "right", "axle": "poros_4", "x": 0.75, "y": 0.25, "mirror_of": "5"},
  {"position": "7", "label": "Tyre 7", "side": "left",  "axle": "poros_4", "x": 0.90, "y": 0.25, "mirror_of": "8"},
  {"position": "8", "label": "Tyre 8", "side": "right", "axle": "poros_4", "x": 0.60, "y": 0.25, "mirror_of": "7"},
  {"position": "9", "label": "Tyre 9", "side": "left",  "axle": "poros_2", "x": 0.17, "y": 0.55, "mirror_of": "10"},
  {"position": "10","label": "Tyre 10","side": "right", "axle": "poros_2", "x": 0.83, "y": 0.55, "mirror_of": "9"},
  {"position": "11","label": "Tyre 11","side": "left",  "axle": "poros_1", "x": 0.35, "y": 0.75, "mirror_of": "12"},
  {"position": "12","label": "Tyre 12","side": "right", "axle": "poros_1", "x": 0.65, "y": 0.75, "mirror_of": "11"}
]'::jsonb
WHERE unit_type = 'SANY_10POS';

-- ── GREADER_6POS ───────────────────────────────────────────────────────────
-- Old: rear/bogie/front at y=0.30/0.55/0.70
-- New: poros_5/poros_2/poros_1
UPDATE unit_type_configs
SET position_config = '[
  {"position": "1", "label": "Tyre 1", "side": "left",  "axle": "poros_5", "x": 0.20, "y": 0.30, "mirror_of": "2"},
  {"position": "2", "label": "Tyre 2", "side": "right", "axle": "poros_5", "x": 0.80, "y": 0.30, "mirror_of": "1"},
  {"position": "3", "label": "Tyre 3", "side": "left",  "axle": "poros_1", "x": 0.30, "y": 0.70, "mirror_of": "4"},
  {"position": "4", "label": "Tyre 4", "side": "right", "axle": "poros_1", "x": 0.70, "y": 0.70, "mirror_of": "3"},
  {"position": "5", "label": "Tyre 5", "side": "left",  "axle": "poros_2", "x": 0.20, "y": 0.55, "mirror_of": "6"},
  {"position": "6", "label": "Tyre 6", "side": "right", "axle": "poros_2", "x": 0.80, "y": 0.55, "mirror_of": "5"}
]'::jsonb
WHERE unit_type = 'GREADER_6POS';

-- ── mounted_position in tyre_master: "Pos N" → "Tyre N" ───────────────────
-- Only updates rows where mounted_position matches the "Pos N" pattern
UPDATE tyre_master
SET mounted_position = 'Tyre ' || substr(mounted_position, 5)
WHERE mounted_position ~ '^Pos [0-9]+$';

-- ── mounted_position in tyre_replacement_details (if exists) ────────────────
-- Covers any historical records that stored position labels
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'tyre_replacement_details'
      AND column_name = 'position'
  ) THEN
    UPDATE tyre_replacement_details
    SET position = 'Tyre ' || substr(position, 5)
    WHERE position ~ '^Pos [0-9]+$';
  END IF;
END $$;

-- ── Generic fallback for any custom templates ───────────────────────────────
-- Covers any unit_type_configs with old-format axle/side keys that don't match
-- the seeded unit types above. Transforms them inline using jsonb_set.
UPDATE unit_type_configs
SET position_config = (
  SELECT jsonb_agg(
    obj ||
    jsonb_build_object(
      'axle',
      (CASE
        WHEN obj->>'axle' = 'front'  THEN 'poros_1'
        WHEN obj->>'axle' = 'bogie'  THEN 'poros_2'
        WHEN obj->>'axle' = 'rear_1' THEN 'poros_3'
        WHEN obj->>'axle' = 'rear_2' THEN 'poros_4'
        WHEN obj->>'axle' = 'rear'   THEN 'poros_5'
        ELSE obj->>'axle'
      END),
      'side',
      (CASE
        WHEN obj->>'side' IN ('front_left','rear_left','rear_1_left','rear_2_left','bogie_left','left')  THEN 'left'
        WHEN obj->>'side' IN ('front_right','rear_right','rear_1_right','rear_2_right','bogie_right','right') THEN 'right'
        ELSE obj->>'side'
      END)
    ) - 'mirror_of'
  )
  FROM jsonb_array_elements(position_config) AS obj
)
WHERE unit_type NOT IN ('ADT_8POS', 'SANY_10POS', 'GREADER_6POS')
  AND (
    position_config::text ~ '"axle"\s*:\s*"(front|bogie|rear_1|rear_2|rear)"'
    OR position_config::text ~ '"side"\s*:\s*"(front_|bogie_|rear_|rear_)"'
  );
