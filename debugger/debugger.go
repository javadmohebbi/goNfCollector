package debugger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Debugger - debugger struct
type Debugger struct {
	debug  bool
	logrus *logrus.Logger
	format string
}

// New - Create new instance of debugger
func New(debug bool, logr *logrus.Logger, format string) *Debugger {
	return &Debugger{
		debug:  debug,
		format: format,
		logrus: logr,
	}
}

// Verbose will verbose log
func (d *Debugger) Verbose(message string, level logrus.Level, flds ...logrus.Fields) {

	// WhereAmI variable to get stack trace
	d.logrus.SetLevel(logrus.DebugLevel)
	w := "false"
	l := strings.ToLower(level.String())

	if d.debug {
		w = d.whereAmI(2)
	}

	// set the log format
	if d.format == "json" {
		d.logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		//LOG
		d.logrus.Formatter = &logrus.TextFormatter{
			// DisableColors:          false,
			DisableLevelTruncation: true,
			ForceColors:            true,
			FullTimestamp:          true,
		}
	}

	// Set fields
	fields := logrus.Fields{
		"debug": w,
	}

	// Set custom fields
	if len(flds) >= 1 {
		fields = flds[0]
	}

	switch l {
	case "debug":
		if d.debug {
			d.logrus.WithFields(fields).Debugln(w, message)
		}
		return
	case "info":
		d.logrus.WithFields(fields).Infoln(message)
		return
	case "warning":
		d.logrus.WithFields(fields).Warningln(message)
		return
	case "fatal":
		d.logrus.WithFields(fields).Fatalln(message)
		return
	case "error":
		d.logrus.WithFields(fields).Errorln(message)
		return
	case "trace":
		d.logrus.WithFields(fields).Traceln(message)
		return
	case "warn":
		d.logrus.WithFields(fields).Warnln(message)
		return
	default:
		return
	}

}

// whereAmI - get the current func and file name
func (d *Debugger) whereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("File: %s Function: %s Line: %d", chopPath(file), runtime.FuncForPC(function).Name(), line)
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}
