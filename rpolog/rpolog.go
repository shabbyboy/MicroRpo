package rpolog

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var (
	logfmt string
	path   string
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	exdir, _ := os.Getwd()
	logdir := filepath.Base(exdir)
	//换成自己的实际路径
	fullpath := "/Users/tugame/newgodemo/microrpo/MicroRpo/runlogs"
	path = filepath.Join(fullpath, logdir)
	//权限需要是0777 否则权限不够
	os.Mkdir(path, 0777)

	logfmt = "/web%v.log.%v"
	log.SetLevel(log.DebugLevel)
	/*
		os.O_WRONLY | os.O_CREATE | O_EXCL        [如果已经存在，则失败】
		os.O_WRONLY | os.O_CREATE                 [如果已经存在，会覆盖写，不会清空原来的文件，而是从头直接覆盖写】
		os.O_WRONLY | os.O_CREATE | os.O_APPEND  【如果已经存在，则在尾部添加写】
	*/
	tmplogfmt := fmt.Sprintf(logfmt, os.Getpid(), time.Now().Format("2006-01-02"))
	fmt.Println(path + tmplogfmt)
	file, err := os.OpenFile(path+tmplogfmt, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	fmt.Println(err)
	log.SetOutput(file)
}

type Level string

const (
	DEBUG Level = "debug"
	INFO  Level = "info"
	WARN  Level = "warn"
	FATAL Level = "fatal"
	ERROR Level = "error"
	PANIC Level = "panic"
)

func logbase() {
	debugogfmt := fmt.Sprintf(logfmt, os.Getpid(), time.Now().Format("2006-01-02"))
	debugpath := path + debugogfmt
	_, err := os.Stat(debugpath)

	if os.IsNotExist(err) {
		file, _ := os.OpenFile(debugpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		log.SetOutput(file)
	}
}

func SetLevel(level Level) {
	switch level {
	case DEBUG:
		log.SetLevel(log.DebugLevel)
	case INFO:
		log.SetLevel(log.InfoLevel)
	case WARN:
		log.SetLevel(log.WarnLevel)
	case FATAL:
		log.SetLevel(log.FatalLevel)
	case ERROR:
		log.SetLevel(log.ErrorLevel)
	case PANIC:
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func Debug(args ...interface{}) {
	logbase()
	log.Debug(args...)
}

func Info(args ...interface{}) {
	logbase()
	log.Info(args...)
}

func Warn(args ...interface{}) {
	logbase()
	log.Warn(args...)
}

func Error(args ...interface{}) {
	logbase()
	log.Error(args...)
}

func Panic(args ...interface{}) {
	logbase()
	log.Panic(args...)
}

func Fatal(args ...interface{}) {
	logbase()
	log.Fatal(args...)
}

func WithFields(fields map[string]interface{}){
	log.WithFields(fields)
}

func WithField(key string,value interface{}){
	log.WithField(key,value)
}
