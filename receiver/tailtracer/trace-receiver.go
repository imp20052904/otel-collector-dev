package tailtracer

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"time"
)

type tailtracerReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (ttr *tailtracerReceiver) Start(ctx context.Context, host component.Host) error {
	ttr.host = host
	ctx = context.Background()
	ctx, ttr.cancel = context.WithCancel(ctx)

	interval, _ := time.ParseDuration(ttr.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ttr.logger.Info("I should start processing traces now!")
				ttr.nextConsumer.ConsumeTraces(ctx, generateTraces(ttr.config.NumberOfTraces))
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (ttr *tailtracerReceiver) Shutdown(ctx context.Context) error {
	ttr.cancel()
	return nil
}
