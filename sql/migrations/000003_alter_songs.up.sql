BEGIN;

ALTER TABLE songs ADD COLUMN group_id INT;

CREATE INDEX IF NOT EXISTS tmp_idx_group_name ON songs (group_name);

INSERT INTO groups (name)
SELECT DISTINCT group_name FROM songs
ON CONFLICT (name) DO NOTHING;

UPDATE songs s
SET group_id = g.id
FROM groups g
WHERE s.group_name = g.name;

ALTER TABLE songs
DROP COLUMN group_name,
ALTER COLUMN group_id SET NOT NULL,
ADD CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(id);

DROP INDEX IF EXISTS tmp_idx_group_name;

COMMIT;
