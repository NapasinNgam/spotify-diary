ALTER TABLE half_year_recaps DROP CONSTRAINT IF EXISTS half_year_recaps_user_id_period_slot_key;
ALTER TABLE half_year_recaps DROP COLUMN IF EXISTS recap_type;
ALTER TABLE half_year_recaps DROP COLUMN IF EXISTS slot_key;
ALTER TABLE half_year_recaps ALTER COLUMN rank SET NOT NULL;
ALTER TABLE half_year_recaps ADD CONSTRAINT half_year_recaps_rank_check CHECK (rank BETWEEN 1 AND 8);
ALTER TABLE half_year_recaps ADD CONSTRAINT half_year_recaps_user_id_period_rank_key UNIQUE(user_id, period, rank);
