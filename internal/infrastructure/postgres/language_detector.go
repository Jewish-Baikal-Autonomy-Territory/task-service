package postgres

import (
	"strings"

	"github.com/pemistahl/lingua-go"
)

type LanguageDetector struct {
	detector lingua.LanguageDetector
}

func (d *LanguageDetector) Detect(text string) string {
	if language, exists := d.detector.DetectLanguageOf(text); exists {
		return strings.ToLower(language.String())
	}
	return "simple"
}

func NewLanguageDetector(accuracy float64) *LanguageDetector {
	internalDetector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(
			lingua.Russian,
			lingua.English,
			lingua.Spanish,
			lingua.French,
			lingua.German,
			lingua.Italian,
			lingua.Portuguese,
			lingua.Dutch,
			lingua.Danish,
			lingua.Finnish,
			lingua.Hungarian,
			lingua.Romanian,
			lingua.Swedish,
			lingua.Turkish).
		WithMinimumRelativeDistance(accuracy).
		Build()
	return &LanguageDetector{
		detector: internalDetector,
	}
}
