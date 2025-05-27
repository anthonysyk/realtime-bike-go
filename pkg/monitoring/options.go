package monitoring

import (
	"errors"
	"go.uber.org/zap"
)

type options struct {
	logger *zap.Logger
}

type Option func(options *options) error

func WithLogger(logger *zap.Logger) Option {
	return func(options *options) error {
		if logger == nil {
			return errors.New("logger is nil")
		}
		options.logger = logger
		return nil
	}
}
