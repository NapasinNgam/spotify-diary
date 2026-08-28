-- Add recap_type column and slot_key for flexible ordering
-- recap_type: 'half_year' (5 ranks + 3 HM) or 'full_year' (10 ranks + 3 HM + golden)
ALTER TABLE half_year_recaps ADD COLUMN IF NOT EXISTS recap_type VARCHAR(20) DEFAULT 'half_year';
ALTER TABLE half_year_recaps ADD COLUMN IF NOT EXISTS slot_key VARCHAR(30) NOT NULL DEFAULT 'rank_1';

-- Remove old rank check constraint and unique
ALTER TABLE half_year_recaps DROP CONSTRAINT IF EXISTS half_year_recaps_rank_check;
ALTER TABLE half_year_recaps ALTER COLUMN rank DROP NOT NULL;

-- Drop old unique and create new one based on slot_key
ALTER TABLE half_year_recaps DROP CONSTRAINT IF EXISTS half_year_recaps_user_id_period_rank_key;
ALTER TABLE half_year_recaps ADD CONSTRAINT half_year_recaps_user_id_period_slot_key UNIQUE(user_id, period, slot_key);

-- Allow empty track_id for text-only slots (overall, summary)
ALTER TABLE half_year_recaps ALTER COLUMN track_id DROP NOT NULL;
