package scheduler

import (
	"context"
	"log/slog"

	"github.com/EmilioQ74/EmBackUp/internal/adapters"
	"github.com/EmilioQ74/EmBackUp/internal/engine"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	s   gocron.Scheduler
	log *slog.Logger
}

func New(log *slog.Logger) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	return &Scheduler{s: s, log: log}, nil
}

func (sc *Scheduler) Add(cron string, eng *engine.Engine, cfg adapters.DBConfig) error {
	_, err := sc.s.NewJob(
		gocron.CronJob(cron, false),
		gocron.NewTask(func() {
			ctx := context.Background()
			if err := eng.Backup(ctx, cfg); err != nil {
				sc.log.Error("scheduled backup failed", "db", cfg.Database, "error", err)
			}
		}),
	)

	return err
}

func (sc *Scheduler) Start() { sc.s.Start() }
func (sc *Scheduler) Stop()  { _ = sc.s.Shutdown() }
