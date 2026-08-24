package ai

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode"

	"github.com/thenumber/app/internal/calc"
)

type GuidanceJSON struct {
	Sleeves   []calc.Sleeve `json:"sleeves"`
	Rationale string        `json:"rationale"`
}

var productTerms = []string{
	"hdfc", "sbi", "icici", "nippon", "uti", "axis", "kotak", "mirae", "motilal",
	"parag", "parikh", "quant", "tata", "aditya", "birla", "dsp", "franklin",
	"invesco", "hsbc", "baroda", "canara", "ppfas", "groww", "zerodha", "coin",
	"nse", "bse", "nifty", "sensex", "amc", "elss", "nfo", "direct plan",
	"regular plan", "ticker", "isin",
}

var productRE = regexp.MustCompile(`(?i)\b(hdfc|sbi|icici|nippon|uti|axis|kotak|mirae|motilal|parag|parikh|quant|tata|aditya|birla|dsp|franklin|invesco|hsbc|baroda|canara|ppfas|groww|zerodha|coin|nifty|sensex|amc)\b`)

var fenceRE = regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)```")

// ParseGuidanceJSON accepts model JSON, maps to allowed categories, and strips product names.
func ParseGuidanceJSON(raw string) (GuidanceJSON, error) {
	raw = strings.TrimSpace(raw)
	if m := fenceRE.FindStringSubmatch(raw); len(m) == 2 {
		raw = strings.TrimSpace(m[1])
	}
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end > start {
		raw = raw[start : end+1]
	}

	var parsed GuidanceJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return GuidanceJSON{}, err
	}

	merged := map[string]float64{}
	order := []string{}
	for _, s := range parsed.Sleeves {
		cat, ok := mapCategory(s.Category)
		if !ok || s.Percent <= 0 {
			continue
		}
		if _, exists := merged[cat]; !exists {
			order = append(order, cat)
		}
		merged[cat] += s.Percent
	}
	if len(merged) == 0 {
		return GuidanceJSON{}, errNoSleeves
	}
	sleeves := make([]calc.Sleeve, 0, len(order))
	for _, cat := range order {
		sleeves = append(sleeves, calc.Sleeve{Category: cat, Percent: merged[cat]})
	}
	sleeves = calc.NormalizePercents(sleeves)
	parsed.Sleeves = sleeves
	parsed.Rationale = stripProducts(parsed.Rationale)
	if strings.TrimSpace(parsed.Rationale) == "" {
		parsed.Rationale = "Category-level mix only. No products named."
	}
	return parsed, nil
}

var errNoSleeves = errString("no usable category sleeves")

type errString string

func (e errString) Error() string { return string(e) }

func mapCategory(name string) (string, bool) {
	cleaned := stripProducts(name)
	cleaned = strings.TrimSpace(cleaned)
	lower := strings.ToLower(cleaned)
	lower = strings.NewReplacer("-", " ", "_", " ").Replace(lower)

	switch {
	case strings.Contains(lower, "international") || strings.Contains(lower, "global") || strings.Contains(lower, "intl"):
		return calc.CatIntl, true
	case strings.Contains(lower, "small"):
		return calc.CatSmall, true
	case strings.Contains(lower, "mid"):
		return calc.CatMid, true
	case strings.Contains(lower, "large") || strings.Contains(lower, "blue chip") || strings.Contains(lower, "bluechip"):
		return calc.CatLarge, true
	case strings.Contains(lower, "gold") || strings.Contains(lower, "sovereign"):
		return calc.CatGold, true
	case strings.Contains(lower, "cash") || strings.Contains(lower, "liquid") || strings.Contains(lower, "money market"):
		return calc.CatCash, true
	case strings.Contains(lower, "debt") || strings.Contains(lower, "bond") || strings.Contains(lower, "gilt") || strings.Contains(lower, "government"):
		return calc.CatDebt, true
	default:
		for _, allowed := range calc.AllowedCategories {
			if strings.EqualFold(cleaned, allowed) {
				return allowed, true
			}
		}
		return "", false
	}
}

func stripProducts(s string) string {
	s = productRE.ReplaceAllString(s, "")
	// Collapse leftover punctuation/spaces.
	var b strings.Builder
	prevSpace := true
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	out := strings.TrimSpace(b.String())
	out = strings.ReplaceAll(out, "  ", " ")
	return out
}

func HasProductName(s string) bool {
	low := strings.ToLower(s)
	for _, t := range productTerms {
		if strings.Contains(low, t) {
			return true
		}
	}
	return false
}
