package logger

import "github.com/sirupsen/logrus"

// utcFormatter for logrus UTC timezone
type utcFormatter struct {
	logrus.Formatter
	service string
	version string
}

// Format logrus timezone as UTC
func (f utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	e.Data["service_name"] = f.service
	e.Data["service_version"] = f.version
	return f.Formatter.Format(e)
}
