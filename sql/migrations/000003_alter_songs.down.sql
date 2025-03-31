ALTER TABLE songs ADD COLUMN group_name TEXT;

UPDATE songs s
SET group_name = g.name
FROM groups g
WHERE s.group_id = g.id;

ALTER TABLE songs
DROP CONSTRAINT IF EXISTS fk_group,
DROP COLUMN IF EXISTS group_id;
