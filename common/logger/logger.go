package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// FuncHandler is a custom slog.Handler that adds the calling function's name
// to the log record. It wraps another slog.Handler.
type FuncHandler struct {
	slog.Handler
}

// Handle intercepts the log record, finds the calling function's name from the
// program counter (PC), adds it as an attribute, and then passes the record
// to the wrapped handler.
func (h *FuncHandler) Handle(ctx context.Context, r slog.Record) error {
	// The slog.Record contains the PC (Program Counter) of the log call site.
	// We use it to get the function, file, and line number.
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		// The f.Function value is a full path like:
		// "github.com/your-module/your-project/main.doSomeWork"
		// We can trim the path to get just "main.doSomeWork".
		funcName := f.Function
		if lastSlash := strings.LastIndex(funcName, "/"); lastSlash != -1 {
			funcName = funcName[lastSlash+1:]
		}

		// Add the function name as an attribute to the log record.
		r.AddAttrs(slog.String("function", funcName))
	}

	// Pass the modified record to the wrapped handler.
	return h.Handler.Handle(ctx, r)
}

// Log is the global, pre-configured logger instance.
var Log *slog.Logger

// The init function runs automatically when the package is first imported.
func Initialize() {
	// Create our custom handler, wrapping the default JSON handler.
	handler := &FuncHandler{
		Handler: slog.NewJSONHandler(os.Stderr, nil),
	}

	// Create a new logger with our custom handler and assign it to the global variable.
	Log = slog.New(handler)
}
