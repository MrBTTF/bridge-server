package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Fields map[string]interface{}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			pc, _, _, _ := runtime.Caller(10)
			f := runtime.FuncForPC(pc)
			file, line := f.FileLine(pc)
			name := getName(f)

			names := []string{name}
			for i := 0; i < 10; i++ {
				if name == "main.main" {
					break
				}
				pc, _, _, ok := runtime.Caller(10 + i)
				if !ok {
					break
				}
				f := runtime.FuncForPC(pc)

				if !strings.Contains(f.Name(), "bridge-server") {
					continue
				}
				names = append(names, getName(f))
			}
			sort.Sort(sort.Reverse(sort.StringSlice(names)))
			return strings.Join(names, "::"), fmt.Sprintf("%s:%d", file, line)
		},
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Warn(args ...interface{}) {
	logrus.SetReportCaller(true)
	logrus.Warn(args...)
	logrus.SetReportCaller(false)
}

func Error(args ...interface{}) {
	logrus.SetReportCaller(true)
	logrus.Error(args...)
	logrus.SetReportCaller(false)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

func WithFields(fields Fields) *logrus.Entry {
	return logrus.WithFields(logrus.Fields(fields))
}

func getName(f *runtime.Func) string {
	name := f.Name()
	index := strings.LastIndex(name, "/")
	if index != -1 {
		name = name[index+1:]
	}
	return name
}
