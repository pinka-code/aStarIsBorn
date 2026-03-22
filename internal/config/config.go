package config

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type StarsRange string

const (
	StarsLow    StarsRange = "low"    // 1-50
	StarsMedium StarsRange = "medium" // 50-200
	StarsHigh   StarsRange = "high"   // 200+
)

type SizeRange string

const (
	SizeSmall  SizeRange = "small"
	SizeMedium SizeRange = "medium"
	SizeLarge  SizeRange = "large"
)

type Criteria struct {
	MinContributors bool
	Stars           StarsRange
	Language        string
	Size            SizeRange
}

func (c Criteria) Pretty() string {
	parts := []string{}
	parts = append(parts, "public repos")

	if c.Language != "" {
		parts = append(parts, fmt.Sprintf("language: %s", c.Language))
	}

	if c.Stars != "" {
		parts = append(parts, fmt.Sprintf("stars: %s", c.Stars))
	}

	if c.Size != "" {
		parts = append(parts, fmt.Sprintf("size: %s", c.Size))
	}

	if c.MinContributors {
		parts = append(parts, "contributors: > threshold")
	}

	return strings.Join(parts, " | ")
}

func ResolveCriteria(rng *rand.Rand) Criteria {
	forcedLang := os.Getenv("LANGUAGE")
	if forcedLang != "" {
		return Criteria{
			MinContributors: true,
			Stars:           StarsMedium,
			Language:        forcedLang,
			Size:            SizeMedium,
		}
	}

	starsOptions := []StarsRange{
		StarsLow,
		StarsMedium,
		StarsHigh,
	}

	sizeOptions := []SizeRange{
		SizeSmall,
		SizeMedium,
		SizeLarge,
	}

	languages := []string{
		"go", "python", "typescript", "rust", "java", "c++",
	}

	return Criteria{
		MinContributors: rng.Intn(2) == 0,
		Stars:           starsOptions[rng.Intn(len(starsOptions))],
		Language:        languages[rng.Intn(len(languages))],
		Size:            sizeOptions[rng.Intn(len(sizeOptions))],
	}
}
