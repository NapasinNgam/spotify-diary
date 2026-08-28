package job

import (
	"fmt"

	"github.com/NapasinNgam/spotify-diary/internal/config"
	"github.com/NapasinNgam/spotify-diary/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron        *cron.Cron
	cfg         *config.Config
	db          *pgxpool.Pool
	logger      *zap.Logger
	userRepo    *repository.UserRepository
	historyRepo *repository.HistoryRepository
}

func NewScheduler(cfg *config.Config, db *pgxpool.Pool, logger *zap.Logger) *Scheduler {
	return &Scheduler{
		cron:        cron.New(),
		cfg:         cfg,
		db:          db,
		logger:      logger,
		userRepo:    repository.NewUserRepository(db),
		historyRepo: repository.NewHistoryRepository(db),
	}
}

func (s *Scheduler) Start() {
	// Sync listening history every N minutes
	syncSpec := fmt.Sprintf("@every %dm", s.cfg.SyncIntervalMinutes)
	s.cron.AddFunc(syncSpec, func() {
		s.logger.Info("Running sync job")
		s.SyncAllUsers()
	})

	// Generate daily summary at 00:30 every day
	s.cron.AddFunc("30 0 * * *", func() {
		s.logger.Info("Running daily summary generation")
		s.GenerateDailySummaries()
	})

	// Generate monthly summary on 1st of each month at 01:00
	s.cron.AddFunc("0 1 1 * *", func() {
		s.logger.Info("Running monthly summary generation")
		s.GenerateMonthlySummaries()
	})

	s.cron.Start()
	s.logger.Info("Scheduler started",
		zap.Int("sync_interval_minutes", s.cfg.SyncIntervalMinutes),
	)
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
