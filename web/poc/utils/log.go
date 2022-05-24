package utils

import (
	"fmt"
	"nacs/common"
	"nacs/utils/logger"
	"nacs/web/poc/errors"
	myerrors "nacs/web/poc/internal/common/errors"
)

var (
	DebugFlag bool
)

/// Info
func InfoF(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	logger.Info(fmt.Sprintf("%v\n", args...))
}

// Error
func ErrorF(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	logger.Error(fmt.Sprintf("%v\n", args...))
}

// PrintError
func ErrorP(err error) {
	// print stack trace if debug
	if common.RunningInfo.PocDebug {
		if DebugFlag {
			switch customErr := errors.Cause(err).(type) {
			case myerrors.CustomError:
				switch customErr.Type {
				// case myerrors.ConvertInterfaceError:
				// case myerrors.CompileError:
				default:
					logger.Error(fmt.Sprintf("%s: %+v", "PocV Error", err))
				}
			default:
				// raw error
				logger.Error(fmt.Sprintf("%s: %+v", "Raw Error", err))
			}

		} else {
			logger.Error(fmt.Sprintf("%v", err))
		}
	}

}

// Warning
func WarningF(format string, args ...interface{}) {
	logger.Warning(fmt.Sprintf(format, args...))
}

func Warning(args ...interface{}) {
	logger.Warning(fmt.Sprintf("%v\n", args...))
}

// Debug
func DebugF(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Debug(args ...interface{}) {
	logger.Debug(fmt.Sprintf("%v\n", args...))
}
