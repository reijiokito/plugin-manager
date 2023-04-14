package log

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

type Log struct {
	Logger hclog.Logger
}

func (l *Log) InitLog(name string) *Log {
	return &Log{
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   name,
			Output: os.Stdout,
			Level:  hclog.Trace,
		}),
	}
}
