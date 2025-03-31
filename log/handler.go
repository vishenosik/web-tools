package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/fatih/color"
	"github.com/vishenosik/web-tools/colors"
	"gopkg.in/yaml.v2"
)

const (
	timeFormat = "15:04:05.000"
)

type Handler struct {
	handler slog.Handler
	writer  io.Writer
	rec     nextFunc
	buf     *bytes.Buffer
	mutex   *sync.Mutex

	outputEmptyAttrs bool
	// syntax highlighter
	highlight *colors.Higlighter

	// marshaller type
	marshalType uint8
}

// The signature of the function for setting parameters
type optsFunc func(*Handler)

func NewHandler(
	writer io.Writer,
	level slog.Level,
	opts ...optsFunc,
) *Handler {

	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	buf := &bytes.Buffer{}

	h := &Handler{
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaultAttrs(handlerOptions.ReplaceAttr),
		}),
		buf:       buf,
		writer:    writer,
		highlight: colors.NewHighlighter(),
		rec:       handlerOptions.ReplaceAttr,
		mutex:     &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {

	attrs, err := h.computeAttrs(ctx, rec)
	if err != nil {
		return err
	}

	attrsStr, err := h.marshal(attrs)
	if err != nil {
		return err
	}

	output := fmt.Sprintf(
		"[%s] %s: %s\n",
		rec.Time.Format(timeFormat),
		level(rec),
		color.CyanString(rec.Message),
	)

	attrsStr = h.highlight.HighlightNumbers(attrsStr)
	attrsStr = h.highlight.HighlightKeyWords(attrsStr)

	if attrsStr != "" {
		output = fmt.Sprintf(
			"%s%s\n",
			output,
			attrsStr,
		)
	}

	_, err = io.WriteString(h.writer, output)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		handler:     h.handler.WithAttrs(attrs),
		writer:      h.writer,
		highlight:   h.highlight,
		buf:         h.buf,
		rec:         h.rec,
		mutex:       h.mutex,
		marshalType: h.marshalType,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		handler:     h.handler.WithGroup(name),
		writer:      h.writer,
		highlight:   h.highlight,
		buf:         h.buf,
		rec:         h.rec,
		mutex:       h.mutex,
		marshalType: h.marshalType,
	}
}

func level(rec slog.Record) string {
	level := rec.Level.String()

	switch rec.Level {

	case slog.LevelDebug:
		level = color.MagentaString(level)

	case slog.LevelInfo:
		level = color.BlueString(level)

	case slog.LevelWarn:
		level = color.YellowString(level)

	case slog.LevelError:
		level = color.RedString(level)

	}
	return level
}

func WithNumbersHighlight(color colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithNumbersHighlight(color))
	}
}

func WithKeyWordsHighlight(keywordsToColors map[string]colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithKeyWordsHighlight(keywordsToColors))
	}
}

const (
	json_marshaller uint8 = iota
	yaml_marshaller
)

type nextFunc func([]string, slog.Attr) slog.Attr

type attrsMap map[string]any

func WithYamlMarshaller() optsFunc {
	return func(h *Handler) {
		h.marshalType = yaml_marshaller
	}
}

func WithJsonMarshaller() optsFunc {
	return func(h *Handler) {
		h.marshalType = json_marshaller
	}
}

func suppressDefaultAttrs(
	next nextFunc,
) nextFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey, slog.LevelKey, slog.MessageKey:
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func (h *Handler) computeAttrs(
	ctx context.Context,
	rec slog.Record,
) (attrsMap, error) {

	h.mutex.Lock()
	defer func() {
		h.buf.Reset()
		h.mutex.Unlock()
	}()

	if err := h.handler.Handle(ctx, rec); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs attrsMap
	err := json.Unmarshal(h.buf.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func (h *Handler) marshal(attrs attrsMap) (string, error) {
	var (
		data []byte
		err  error
	)
	if len(attrs) > 0 {

		switch h.marshalType {
		case json_marshaller:
			data, err = json.MarshalIndent(attrs, "", "  ")

		case yaml_marshaller:
			data, err = yaml.Marshal(attrs)
		}

		if err != nil {
			return "", err
		}
	}
	return string(data), nil
}
