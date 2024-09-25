package pkg

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}

// Debug выводит сообщение на уровне debug.
func Debug(message string, fields map[string]interface{}) {
	Log.WithFields(fields).Debug(message)
}

// Info выводит сообщение на уровне info.
func Info(message string, fields map[string]interface{}) {
	Log.WithFields(fields).Info(message)
}

// Error выводит сообщение на уровне error.
func Error(message string, fields map[string]interface{}) {
	Log.WithFields(fields).Error(message)
}

// Fatal выводит сообщение на уровне fatal и завершает программу.
func Fatal(message string, fields map[string]interface{}) {
	Log.WithFields(fields).Fatal(message)
}
