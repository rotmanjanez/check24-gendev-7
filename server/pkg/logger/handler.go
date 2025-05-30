package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// TextHandler is a custom text handler for slog
type TextHandler struct {
	w           io.Writer
	opts        slog.HandlerOptions
	mu          sync.Mutex
	attrsPrefix string
	attrs       []slog.Attr
}

// NewTextHandler creates a new TextHandler
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) *TextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &TextHandler{
		w:           w,
		opts:        *opts,
		attrsPrefix: "",
		attrs:       []slog.Attr{},
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *TextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle handles the Record.
func (h *TextHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String()
	timeStr := r.Time.Format(time.RFC3339)

	// Find and extract provider attribute
	var provider string = "root"
	var attrs []string

	// First add handler's attrs
	for _, attr := range h.attrs {
		if attr.Key == "provider" || attr.Key == "adapter" {
			provider = attr.Value.String()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		}
	}

	// Then add record's attrs
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "provider" || attr.Key == "adapter" {
			provider = attr.Value.String()
			// Don't add provider to attrs slice since we'll handle it specially
			return true
		}
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		return true
	})

	// Format the log message
	h.mu.Lock()
	defer h.mu.Unlock()

	// Colorize level based on severity
	var levelColored string
	switch level {
	case "ERROR":
		levelColored = "\033[31m" + level + "\033[0m" // Red for ERROR
	case "WARN":
		levelColored = "\033[33m" + level + "\033[0m" // Yellow for WARN
	case "INFO":
		levelColored = "\033[32m" + level + "\033[0m" // Green for INFO
	case "DEBUG":
		levelColored = "\033[36m" + level + "\033[0m" // Cyan for DEBUG
	default:
		levelColored = level
	}

	// Colorize provider with a unique color
	providerColored := fmt.Sprintf("\033[35m%s\033[0m", provider) // Magenta for provider

	// Format and write the log message
	message := fmt.Sprintf(
		"%s %s [%s] %s %s\n",
		timeStr,
		levelColored,
		providerColored,
		r.Message,
		strings.Join(attrs, " "),
	)

	_, err := h.w.Write([]byte(message))
	return err
}

// WithAttrs returns a new handler with the given attributes.
func (h *TextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler and copy all attributes
	newHandler := &TextHandler{
		w:           h.w,
		opts:        h.opts,
		attrsPrefix: h.attrsPrefix,
		attrs:       append(h.attrs, attrs...),
	}
	return newHandler
}

// WithGroup returns a new handler with the given group name.
func (h *TextHandler) WithGroup(name string) slog.Handler {
	// Create a new handler with the group name as prefix
	newHandler := &TextHandler{
		w:           h.w,
		opts:        h.opts,
		attrsPrefix: h.attrsPrefix + name + ".",
		attrs:       h.attrs,
	}
	return newHandler
}
