package log

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the debug severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}

type Logger struct {
	*logrus.Entry
}

// Log global logger
var Log Logger = Wrap(logrus.WithFields(logrus.Fields{}))

// LoggerWithSpan spwan new logger with trace span
func (l Logger) LoggerWithSpan(ctx context.Context) Logger {
	span := trace.FromContext(ctx)
	if span == nil {
		return l
	}

	spanID := span.SpanContext().SpanID.String()
	traceID := span.SpanContext().TraceID.String()

	newLogrus := l.WithFields(logrus.Fields{
		"trace": traceID,
		"span":  spanID,
	})

	return Wrap(newLogrus)
}

// Wrap wrap *logrus.Entry to log.Logger
func Wrap(logrusEntry *logrus.Entry) Logger {
	return Logger{logrusEntry}
}
