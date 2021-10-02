package db

import (
	"context"
	"fmt"
	"time"

	"github.com/phuslu/log"
	"gorm.io/gorm/logger"
)

type Logger struct {
	Log    *log.Logger
	Config logger.Config
	Slient bool
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	var newLogger = Logger{Log: l.Log}
	switch level {
	case logger.Silent:
		newLogger.Slient = true
	case logger.Error:
		newLogger.Log.SetLevel(log.ErrorLevel)
	case logger.Warn:
		newLogger.Log.SetLevel(log.WarnLevel)
	case logger.Info:
		newLogger.Log.SetLevel(log.InfoLevel)
	}

	return &newLogger
}

func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	l.Log.Info().Msgf(format, args...)
}

func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.Log.Warn().Msgf(format, args...)
}

func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	l.Log.Error().Msgf(format, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.Slient {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.Log.Level >= log.ErrorLevel:
		sql, rows := fc()
		if rows == -1 {
			l.Log.Error().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Msg("")
		} else {
			l.Log.Error().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Int64("rows", rows).Msg("")
		}
	case elapsed > l.Config.SlowThreshold && l.Config.SlowThreshold != 0 && l.Log.Level >= log.WarnLevel:
		sql, rows := fc()
		if rows == -1 {
			l.Log.Warn().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Msgf("SLOW SQL >= %v", l.Config.SlowThreshold)
		} else {
			l.Log.Warn().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Int64("rows", rows).Msgf("SLOW SQL >= %v", l.Config.SlowThreshold)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.Log.Info().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Msg("")
		} else {
			l.Log.Info().Caller(10, false).Err(err).Str("elapsed", fmt.Sprintf("%s", elapsed)).Str("sql", sql).Int64("rows", rows).Msg("")
		}
	}
}
