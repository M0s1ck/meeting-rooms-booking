package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sched  gocron.Scheduler
	logger *slog.Logger
}

func New(logger *slog.Logger) (*Scheduler, error) {
	s, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	if err != nil {
		return nil, err
	}

	return &Scheduler{sched: s, logger: logger}, nil
}

func (s *Scheduler) AddJob(
	ctx context.Context,
	interval time.Duration,
	name string,
	job func(context.Context) error,
) error {

	_, err := s.sched.NewJob(
		gocron.DurationJob(interval),
		gocron.NewTask(func() {
			s.runJob(ctx, name, job)
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	return err
}

func (s *Scheduler) AddDailyJobAt(
	ctx context.Context,
	hour, minute uint,
	name string,
	job func(context.Context) error,
) error {
	_, err := s.sched.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(hour, minute, 0))),
		gocron.NewTask(func() {
			s.runJob(ctx, name, job)
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	return err
}

func (s *Scheduler) Start() {
	s.sched.Start()
}

func (s *Scheduler) Shutdown() error {
	return s.sched.Shutdown()
}

// logs start, finish, time executing
func (s *Scheduler) runJob(ctx context.Context, name string, job func(context.Context) error) {
	start := time.Now()
	s.logger.Info("job started", "job", name)

	err := job(ctx)

	duration := time.Since(start)

	if err != nil {
		s.logger.Error("job failed",
			"job", name,
			"duration", duration,
			"error", err,
		)
		return
	}

	s.logger.Info("job finished", "job", name, "duration", duration)
}
