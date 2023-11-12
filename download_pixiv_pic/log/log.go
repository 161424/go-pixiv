package log

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"time"
)

type Tp int

const (
	LogFiles Tp = 1 << iota
	LogStdouts
)

var logs = &Logs{}

// type L = logs
type Logs struct {
	logFile   *slog.Logger
	logStdout *slog.Logger
}

var opt slog.HandlerOptions
var logFile *os.File

func init() {

	now := time.Now()
	// 创建log文件
	var err error
	if _, err = os.Stat("../log/log"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("../log/log", 777)
			if err != nil {
				fmt.Println(2, err)
			}
		}
	}
	// 创建xx.log
	logFile, err = os.OpenFile(fmt.Sprintf("../log/log/%s.log", now.Format("2006-01-02")), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	opt = slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}

			if a.Key == slog.SourceKey {
				_, file, lineNo, ok := runtime.Caller(7)
				if !ok {
				}
				//funcName := runtime.FuncForPC(pc).Name()
				fileName := path.Base(file) // Base函数返回路径的最后一个元素
				//f := runtime.CallersFrames()
				//fmt.Printf(file, lineNo)
				f := fmt.Sprintf("%s:%d", fileName, lineNo)
				a.Value = slog.StringValue(f)

			}
			return a
		},
	}

}

func NewSlogGroup(group string) *Logs {
	logs.logFile = slog.New(slog.NewTextHandler(logFile, &opt).WithGroup(group))
	logs.logStdout = slog.New(slog.NewTextHandler(os.Stdout, &opt).WithGroup(group))

	return logs
}
func Gt() {

}
func NewDefaultSlog() *Logs {
	logs.logFile = slog.New(slog.NewTextHandler(logFile, &opt))
	logs.logStdout = slog.New(slog.NewTextHandler(os.Stdout, &opt))
	return logs
	//slog.LogAttrs()
}

func (l *Logs) Send(level slog.Level, msg string, ls Tp, args ...any) {
	if ls <= 0 {
		return
	}
	if len(args) == 0 {
		args = nil
	}
	lchan := make(chan *slog.Logger, 2)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	for ls > 0 {
		var w Tp = 0
		switch {
		case ls&LogFiles > 0:
			lchan <- l.logFile
			w = LogFiles
		case ls&LogStdouts > 0:
			lchan <- l.logStdout
			w = LogStdouts
		}
		ls &^= w

	}

	close(lchan)
	switch level {
	case slog.LevelInfo:
		for s := range lchan {
			if len(args) == 0 {
				s.Info(msg)
			} else {
				s.Info(msg, args)
			}
		}
	case slog.LevelWarn:
		for s := range lchan {
			if len(args) == 0 {
				s.Warn(msg)
			} else {
				s.Warn(msg, args)
			}
		}

	case slog.LevelError:
		for s := range lchan {
			if len(args) == 0 {
				s.Error(msg)
			} else {
				s.Error(msg, args)
			}
		}

	case slog.LevelDebug:
		for s := range lchan {
			if len(args) == 0 {
				s.Debug(msg)
			} else {
				s.Debug(msg, args)
			}
		}

	default:
		l.Send(slog.LevelWarn, "未知错误", LogFiles|LogStdouts, args)
	}
}

//func (l *Logs) Info(msg string, ls ...*slog.Logger) {
//	for _, s := range ls {
//		s.Info(msg, l.Args)
//	}
//	l.Args = ""
//}
//
//func (l *Logs) Warn(msg string, ls ...*slog.Logger) {
//	for _, s := range ls {
//		s.Warn(msg, l.Args)
//	}
//	l.Args = ""
//}
//
//func (l *Logs) Error(msg string, ls ...*slog.Logger) {
//	for _, s := range ls {
//		s.Error(msg, l.Args)
//	}
//	l.Args = ""
//}
//
//func (l *Logs) Debug(msg string, ls ...*slog.Logger) {
//	for _, s := range ls {
//		s.Debug(msg, l.Args)
//	}
//	l.Args = ""
//
//}
