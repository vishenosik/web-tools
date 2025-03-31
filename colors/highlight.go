package colors

import (
	// builtin
	"regexp"

	// internal
	"github.com/vishenosik/web-tools/regex"
)

type Higlighter struct {
	numbers          *regexp.Regexp
	doNumbers        bool
	numbersColor     ColorCode
	keywords         *regexp.Regexp
	keywordsToColors map[string]ColorCode
}

// optsFunc is a function type used for configuring a Higlighter instance.
//
// It takes a pointer to a Higlighter as its parameter and returns nothing.
// This type is typically used in a functional options pattern to provide
// a flexible way of configuring Higlighter objects.
//
// Parameters:
//   - h: A pointer to the Higlighter instance to be configured.
type optsFunc func(*Higlighter)

// NewHighlighter creates and returns a new Higlighter instance.
//
// It initializes a Higlighter with default settings and applies any provided
// option functions to customize the highlighter.
//
// Parameters:
//   - opts: A variadic parameter of option functions (optsFunc) that can be used
//     to configure the Higlighter. Each function in opts will be called with the
//     Higlighter instance, allowing for custom configuration.
//
// Returns:
//   - *Higlighter: A pointer to the newly created and configured Higlighter instance.
func NewHighlighter(
	opts ...optsFunc,
) *Higlighter {
	h := &Higlighter{
		numbers: regex.NumberRegex,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Modify updates an existing Higlighter or creates a new one if the input is nil.
// It applies the provided option functions to configure the Higlighter.
//
// Parameters:
//   - h: A pointer to the Higlighter to be modified. If nil, a new Higlighter is created.
//   - opts: A variadic parameter of option functions (optsFunc) used to configure the Higlighter.
//
// Returns:
//   - *Higlighter: A pointer to the modified or newly created Higlighter instance.
func Modify(
	h *Higlighter,
	opts ...optsFunc,
) *Higlighter {
	if h == nil {
		h = NewHighlighter(opts...)
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HighlightNumbers applies color highlighting to numbers in the input string if enabled.
//
// Parameters:
//   - src: The input string to be processed.
//
// Returns:
//   - string: The processed string with numbers highlighted if enabled, otherwise the original string.
func (h *Higlighter) HighlightNumbers(src string) string {
	if h.doNumbers {
		return h.numbers.ReplaceAllStringFunc(src, func(s string) string { return colors[h.numbersColor](s) })
	}
	return src
}

// HighlightKeyWords applies color highlighting to specified keywords in the input string.
//
// Parameters:
//   - src: The input string to be processed.
//
// Returns:
//   - string: The processed string with keywords highlighted if any are defined, otherwise the original string.
func (h *Higlighter) HighlightKeyWords(src string) string {
	if len(h.keywordsToColors) == 0 {
		return src
	}
	return h.keywords.ReplaceAllStringFunc(src, func(s string) string {
		return colors[h.keywordsToColors[s]](s)
	})
}

// WithNumbersHighlight returns an option function that enables number highlighting with the specified color.
//
// Parameters:
//   - color: The ColorCode to be used for highlighting numbers.
//
// Returns:
//   - optsFunc: A function that configures a Higlighter to highlight numbers with the specified color.
func WithNumbersHighlight(color ColorCode) optsFunc {
	return func(h *Higlighter) {
		h.doNumbers = true
		h.numbersColor = color
	}
}

// WithKeyWordsHighlight returns an option function that enables keyword highlighting with specified colors.
//
// Parameters:
//   - keywordsToColors: A map where keys are keywords to highlight and values are their corresponding ColorCodes.
//
// Returns:
//   - optsFunc: A function that configures a Higlighter to highlight keywords with their specified colors.
func WithKeyWordsHighlight(keywordsToColors map[string]ColorCode) optsFunc {
	return func(h *Higlighter) {
		keywords := make([]string, 0, len(keywordsToColors))
		for key := range keywordsToColors {
			keywords = append(keywords, key)
		}
		h.keywords = regex.KeyWordsCompile(keywords...)
		h.keywordsToColors = keywordsToColors
	}
}
