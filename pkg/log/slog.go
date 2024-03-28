package log

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}
type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}
	timeStr := r.Time.Format("[2006-01-02 15:04:05]")
	msg := color.CyanString(r.Message)
	filePath := getPrefix()
	h.l.Println(timeStr, filePath, level, msg, color.WhiteString(string(b)))
	return nil
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}

func selectLevel(s string) (l slog.Leveler) {
	level := strings.ToLower(s)
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelDebug
	}
	return
}

func getCallTrace() (string, int) {
	_, file, lineNo, ok := runtime.Caller(5)
	if ok {
		return file, lineNo
	} else {
		return "", 0
	}
}

func getPrefix() string {
	file, lineNo := getCallTrace()
	path := strings.Split(file, "/")
	if len(path) > 3 {
		file = strings.Join(path[len(path)-1:], "/")
	}
	return file + ":" + strconv.Itoa(lineNo) + " "
}

func NewSlog(conf *viper.Viper) *slog.Logger {
	level := selectLevel(conf.GetString("log.level"))
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: level,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	return logger
}
