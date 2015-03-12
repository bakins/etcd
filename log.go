package main

import (
	"bytes"
	"time"

	"github.com/Sirupsen/logrus"
)

// LogFortmatter is a simple formatter for logrus
type LogFormatter struct {
	DisableTimestamps bool
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if !f.DisableTimestamps {
		b.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	}
	b.WriteString(entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}
