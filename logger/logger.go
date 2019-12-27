package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

//Environment the type for which the logger is c reated
type Environment string

const (
	//PRODUCTION the production environment
	PRODUCTION Environment = "production"

	//DEVELOPMENT the development environment
	DEVELOPMENT Environment = "development"

	//TESTING the testing environment
	TESTING Environment = "testing"
)

//IsEqual checks if the given env is is equal to this environment
func (env Environment) IsEqual(compEnv Environment) bool {
	return env == compEnv
}

//GetLoggerEnvironment returns the logger environment as per the given string
func GetLoggerEnvironment(env string) Environment {
	switch env {
	case "development":
		return DEVELOPMENT
	case "production":
		return PRODUCTION
	case "testing":
		return TESTING
	default:
		return ""
	}
}

//RecordTenure the tenure for which the details should be logged in a single file
type RecordTenure int

const (
	//Daily the logger will create a new file for logs everyday
	Daily RecordTenure = 0

	//Weekly the logger will create a new file for logs every week
	Weekly RecordTenure = 1

	//Monthly the logger will create a new file for logs every month
	Monthly RecordTenure = 2

	//Yearly the logger will create a new file for logs every year
	Yearly RecordTenure = 3

	//Forever all logs will be stored in one single file
	Forever RecordTenure = 4
)

//LogWriter writer context for every log
type LogWriter interface {
	WriteLog(string)
}

//LogReader reader context to file log
type LogReader interface {
	GetLog() string
}

//FormatTimeFunc the function that should be implemented to format the time string
type FormatTimeFunc func(time.Time) string

//Logger the logger instance
type Logger struct {
	logPath        string
	fileNameSuffix string
	formatTime     FormatTimeFunc
	environment    Environment
	tenure         RecordTenure
}

//NewLogger creates and returns a new Logger instance
func NewLogger(path string, fileNameSuffix string, fTime FormatTimeFunc, env Environment, tenure RecordTenure) *Logger {
	logger := &Logger{
		logPath:        path,
		fileNameSuffix: fileNameSuffix,
		formatTime:     fTime,
		environment:    env,
		tenure:         tenure,
	}
	return logger
}

func (logger Logger) getTime(time time.Time) string {
	if logger.formatTime != nil {
		return logger.formatTime(time)
	}
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second())
}

//LogWithReader Dumps the final log available from the reader into the output source
func (logger *Logger) LogWithReader(reader LogReader) {
	logger.log(reader.GetLog())
}

func (logger *Logger) log(message string) {
	if logger.logPath != "" {
		f, err := os.OpenFile(logger.getFileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			defer f.Close()
			if _, err = f.WriteString(message); err == nil {
				return
			}
		}
		fmt.Println("Unable to log in log file ", err.Error())
		fmt.Println(message)
	}
	fmt.Println(message)
}

func (logger *Logger) getFileName() string {
	fileName := ""
	currTime := time.Now()
	switch logger.tenure {
	case Daily:
		year, month, date := currTime.Date()
		fileName = fmt.Sprint(year, "-", month, "-", date, "-", logger.fileNameSuffix, ".txt")
		return path.Join(logger.logPath, fileName)
	case Weekly:
		year, week := currTime.ISOWeek()
		fileName = fmt.Sprint(year, "-", week, "-", logger.fileNameSuffix, ".txt")
		return path.Join(logger.logPath, fileName)
	case Monthly:
		year, month, _ := currTime.Date()
		fileName = fmt.Sprint(year, "-", month, "-", logger.fileNameSuffix, ".txt")
		return path.Join(logger.logPath, fileName)
	case Yearly:
		year, _, _ := currTime.Date()
		fileName = fmt.Sprint(year, "-", logger.fileNameSuffix, ".txt")
		return path.Join(logger.logPath, fileName)
	case Forever:
		fileName = fmt.Sprint(logger.fileNameSuffix, ".txt")
		return path.Join(logger.logPath, fileName)
	}

	return ""
}

//Log logs the given data into the writer context
func (logger *Logger) Log(level Level, message string, trace string, time time.Time, fields map[string]interface{}, writer LogWriter) {

	var field strings.Builder
	if fields != nil {
		for key, val := range fields {
			field.WriteString(fmt.Sprint(key, " ", val))
		}
	}
	_, file, line, _ := runtime.Caller(2)
	if writer == nil {
		logger.log(fmt.Sprint(level.String(), "   ", logger.getTime(time), "   ", file, ":", line, "	", message, "   ", trace, "   ", field.String(), "\n"))
		return
	}
	writer.WriteLog(fmt.Sprint(level.String(), "   ", logger.getTime(time), "   ", file, ":", line, "	", message, "   ", trace, "   ", field.String(), "\n"))
}

//LogInfo logs the info in the writer context
func (logger *Logger) LogInfo(message string, fields map[string]interface{}, writer LogWriter) {
	logger.Log(Info, message, "", time.Now(), fields, writer)
}

// //LogWarning logs the warning in the writer context
// func (logger Logger) LogWarning(message string, fields map[string]interface{}, writer LogWriter) {
// 	logger.Log(Warning, message, "", time.Now(), fields, writer)
// }

//LogSevere logs the severe log in the writer context
func (logger *Logger) LogSevere(message string, fields map[string]interface{}, writer LogWriter) {
	logger.Log(Severe, message, "", time.Now(), fields, writer)
}

// //LogFatal logs the fatal log in the writer context
// func (logger Logger) LogFatal(message string, fields map[string]interface{}, writer LogWriter) {
// 	logger.Log(Fatal, message, string(debug.Stack()), time.Now(), fields, writer)
// }

//LogPanic logs the panic log in the writer context
func (logger Logger) LogPanic(message string, fields map[string]interface{}, writer LogWriter) {
	logger.Log(Panic, message, string(debug.Stack()), time.Now(), fields, writer)
}

// //LogDebug logs the debug log in the writer context, will log for every environment other than the production environment
// func (logger Logger) LogDebug(message string, fields map[string]interface{}, writer LogWriter) {
// 	if logger.environment != PRODUCTION {
// 		logger.Log(Debug, message, "", time.Now(), fields, writer)
// 	}
// }

// //LogTrace logs the trace log in the writer context
// func (logger Logger) LogTrace(writer LogWriter) {
// 	logger.Log(Trace, "", string(debug.Stack()), time.Now(), nil, writer)
// }
