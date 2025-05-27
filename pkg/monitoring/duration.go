package monitoring

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"runtime"
	"time"
)

func MeasureExecutionTimeWithCtxError(ctx context.Context, f func(ctx context.Context) error, opts ...Option) (time.Duration, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return time.Duration(0), err
		}
	}
	start := time.Now()
	err := f(ctx)
	duration := time.Since(start)
	if options.logger != nil {
		funcInfo := runtime.FuncForPC(reflect.ValueOf(f).Pointer())
		options.logger.Info(fmt.Sprintf("%s execution time", funcInfo.Name()), zap.String("duration", duration.String()))
	}
	return duration, err
}
