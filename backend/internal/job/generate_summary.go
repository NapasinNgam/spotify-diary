package job

import (
	"go.uber.org/zap"
)

// GenerateDailySummaries creates daily summaries for yesterday (all users)
func (s *Scheduler) GenerateDailySummaries() {
	// TODO: Query listening_history for yesterday
	// GROUP BY user_id → compute stats → INSERT INTO daily_summaries
	s.logger.Info("Daily summary generation - TODO: implement")
}

// GenerateMonthlySummaries creates monthly summaries for the previous month (all users)
func (s *Scheduler) GenerateMonthlySummaries() {
	// TODO: Query listening_history for previous month
	// Compute top 10 tracks, top artists, genre breakdown
	// INSERT INTO monthly_summaries
	s.logger.Info("Monthly summary generation - TODO: implement",
		zap.String("note", "Should run on 1st of each month for previous month"),
	)
}
