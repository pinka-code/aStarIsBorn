package config

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