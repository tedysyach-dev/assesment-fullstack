package config

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type slowQueryHook struct {
	logger    *logrus.Logger
	threshold time.Duration
}

func newSlowQueryHook(logger *logrus.Logger, threshold time.Duration) *slowQueryHook {
	return &slowQueryHook{logger: logger, threshold: threshold}
}

func (h *slowQueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *slowQueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	duration := time.Since(event.StartTime)
	if duration >= h.threshold {
		h.logger.WithFields(logrus.Fields{
			"duration": duration,
			"query":    event.Query,
		}).Warn("Slow query detected")
	}
}
