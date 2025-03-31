CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_songs_group_id ON songs (group_id);
CREATE INDEX IF NOT EXISTS idx_songs_release_date ON songs (release_date);
CREATE INDEX IF NOT EXISTS idx_songs_text_search ON songs USING GIN (text gin_trgm_ops);
