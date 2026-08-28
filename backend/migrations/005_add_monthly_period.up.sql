-- Add period_start, period_end, source to monthly_summaries
ALTER TABLE monthly_summaries ADD COLUMN IF NOT EXISTS period_start DATE;
ALTER TABLE monthly_summaries ADD COLUMN IF NOT EXISTS period_end DATE;
ALTER TABLE monthly_summaries ADD COLUMN IF NOT EXISTS source VARCHAR(30) DEFAULT 'listening_history';
