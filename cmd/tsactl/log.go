package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

//TODO: refactor

type CustomHandler struct{}

func (h *CustomHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("2006-01-02T15:04:05.000Z07:00")
	level := fmt.Sprintf("%-5s", r.Level.String())
	msg := r.Message

	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		val := a.Value.Any()

		switch v := val.(type) {
		case int, int64, int32, uint, uint64, float64, float32, bool:
			attrs = append(attrs, fmt.Sprintf("%s=%v", key, v))

		case []byte:
			if isLikelyBinary(v) {
				attrs = append(attrs, fmt.Sprintf("%s=0x%s", key, hex.EncodeToString(v)))
			} else {
				s := string(v)
				s = escapeString(s)
				attrs = append(attrs, fmt.Sprintf("%s=\"%s\"", key, s))
			}

		case string:
			if isLikelyBinary([]byte(v)) {
				attrs = append(attrs, fmt.Sprintf("%s=0x%s", key, hex.EncodeToString([]byte(v))))
			} else {
				attrs = append(attrs, fmt.Sprintf("%s=\"%s\"", key, escapeString(v)))
			}

		default:
			s := fmt.Sprintf("%v", v)
			s = escapeString(s)
			attrs = append(attrs, fmt.Sprintf("%s=\"%s\"", key, s))
		}
		return true
	})

	attrStr := ""
	if len(attrs) > 0 {
		attrStr = ", " + strings.Join(attrs, ", ")
	}

	_, _ = fmt.Fprintf(os.Stderr, "%s [%s] %s%s\n", timestamp, level, msg, attrStr)
	return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *CustomHandler) WithGroup(name string) slog.Handler       { return h }

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}

func isLikelyBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	nonPrintable := 0
	for _, b := range data {
		if (b < 0x20 || b > 0x7E) && b != '\n' && b != '\r' && b != '\t' {
			nonPrintable++
		}
	}
	return float64(nonPrintable)/float64(len(data)) > 0.3
}
