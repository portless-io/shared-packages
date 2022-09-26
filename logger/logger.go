package logger

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warningf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	Debugw(text string, context map[string]interface{})
	Infow(text string, context map[string]interface{})
	Warningw(text string, context map[string]interface{})
	Errorw(text string, context map[string]interface{})
	Fatalw(text string, context map[string]interface{})
}
