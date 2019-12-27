package logger

//Level type of level for logging
type Level int

const (
	//Info level info
	Info Level = 0

	//Warning level warning
	Warning Level = 1

	//Severe level severe
	Severe Level = 2

	//Fatal level fatal
	Fatal Level = 3

	//Panic level panic
	Panic Level = 4

	//Trace level trace
	Trace Level = 5

	//Debug level debug
	Debug Level = 6
)

var levels = [7]string{"INFO", "WARNING", "SEVERE", "FATAL", "PANIC", "TRACE", "DEBUG"}

//String return the printable value of the Level
func (level Level) String() string {
	if isValidLevel(level) {
		return levels[level]
	}
	return ""
}

//method to check if the given level is valid
func isValidLevel(level Level) bool {
	if level < 0 || level > 6 {
		return false
	}
	return true
}
