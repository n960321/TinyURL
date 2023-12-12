package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// LogLevel log level
type LogLevel int

type gormLogger struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func New(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &gormLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Info().Msgf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Warn().Msgf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Error().Msgf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			log.Error().Msgf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Error().Msgf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.Warn().Msgf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warn().Msgf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Info().Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Info().Msgf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
