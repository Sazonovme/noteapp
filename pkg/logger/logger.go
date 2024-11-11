package logger

import "github.com/sirupsen/logrus"

// PanicLevel 0 ignor
// FatalLevel 1 <- Дальнешая работа невозможна
// ErrorLevel 2 <- Работа продоолжается, но требуются срочные доработки
// WarnLevel 3 <- Обратить внимание
// InfoLevel 4 ignor
// DebugLevel 5 <- // Информационный уровень (отработка каких то функций - данные)

func NewLog(packageFuncName string, logLevel int, err error, errMessage string, data interface{}) {
	switch logLevel {
	case 5:
		logrus.WithFields(logrus.Fields{
			"place": packageFuncName,
			"err":   err,
			"data":  data,
		}).Debug(errMessage)
	case 3:
		logrus.WithFields(logrus.Fields{
			"place": packageFuncName,
			"err":   err,
			"data":  data,
		}).Warn(errMessage)
	case 2:
		logrus.WithFields(logrus.Fields{
			"place": packageFuncName,
			"err":   err,
			"data":  data,
		}).Error(errMessage)
	case 1:
		logrus.WithFields(logrus.Fields{
			"place": packageFuncName,
			"err":   err,
			"data":  data,
		}).Fatal(errMessage)
	}
}
