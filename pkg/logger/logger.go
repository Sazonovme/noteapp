package logger

import "github.com/sirupsen/logrus"

// PanicLevel 0 ignor
// FatalLevel 1 <-
// ErrorLevel 2 <-
// WarnLevel 3 <-
// InfoLevel 4 ignor
// DebugLevel 5 <-
// TraceLevel 6 <-

func NewLog(packageName string, funcName string, err error, data interface{}, logLevel int, errMessage string) {
	switch logLevel {
	case 6:
		logrus.WithFields(logrus.Fields{
			"package": packageName,
			"func":    funcName,
			"err":     err,
			"data":    data,
		}).Trace(errMessage)
	case 5:
		logrus.WithFields(logrus.Fields{
			"package": packageName,
			"func":    funcName,
			"err":     err,
			"data":    data,
		}).Debug(errMessage)
	case 3:
		logrus.WithFields(logrus.Fields{
			"package": packageName,
			"func":    funcName,
			"err":     err,
			"data":    data,
		}).Warn(errMessage)
	case 2:
		logrus.WithFields(logrus.Fields{
			"package": packageName,
			"func":    funcName,
			"err":     err,
			"data":    data,
		}).Error(errMessage)
	case 1:
		logrus.WithFields(logrus.Fields{
			"package": packageName,
			"func":    funcName,
			"err":     err,
			"data":    data,
		}).Fatal(errMessage)
	}
}
