package verbyndich

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/rotmanjanez/check24-gendev-7/internal/units"
	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

// Common errors that can be returned by the parser
var (
	ErrMultipleMatches    = errors.New("multiple matches found for pattern")
	ErrInvalidMatchFormat = errors.New("invalid match format")
	ErrIncompleteMatch    = errors.New("could not fully match description data")
)

type DescriptionParser struct {
	emptyPattern *regexp.Regexp
	patterns     map[*regexp.Regexp]HandlerFunc
	logger       *slog.Logger
}

type HandlerFunc func(d *DescriptionParser, match []int, description string, target *DescriptionData) error

func NewDescriptionParser(logger *slog.Logger) *DescriptionParser {
	patters := map[string]HandlerFunc{
		`Für nur (\d+)€ im Monat erhalten Sie eine ([a-zA-Z]+)-Verbindung mit einer Geschwindigkeit von (\d+) ([a-zA-Z/]+)\.`:                                    (*DescriptionParser).handlePriceTypeAndSpeed,
		`Bitte beachten Sie, dass die Mindestvertragslaufzeit (\d+) ([a-zA-Z]+) beträgt\.`:                                                                       (*DescriptionParser).handleMinimalContractDuration,
		`Mit diesem Angebot erhalten Sie einen Rabatt von (\d+)% auf Ihre monatliche Rechnung bis zum (\d+)\. Monat\.\s*(Der maximale Rabatt beträgt (\d+)€\.)?`: (*DescriptionParser).handlePercentageDiscount,
		`Ab dem (\d+)\. Monat beträgt der monatliche Preis (\d+)€\.`:                                                                                             (*DescriptionParser).handleLongRunningPrice,
		`Ab (\d+)([a-zA-Z]+) pro Monat wird die Geschwindigkeit gedrosselt\.`:                                                                                    (*DescriptionParser).handleUnthrottledCapacity,
		`Zögern Sie nicht und schlagen Sie jetzt zu\!`:                                                                                                           (*DescriptionParser).noOp,
		`Dieses einzigartige Angebot ist der perfekte Match für Sie\.`:                                                                                           (*DescriptionParser).noOp,
		`Zusätzlich sind folgende Fernsehsender enthalten ([\w\+]+)\.`:                                                                                           (*DescriptionParser).handleAdditionalTVChannels,
		`Dieses Angebot ist nur für Personen unter (\d+) Jahren verfügbar\.`:                                                                                     (*DescriptionParser).handleMaxAge,
		`Dieses Angebot ist nur für Personen über (\d+) Jahren verfügbar\.`:                                                                                      (*DescriptionParser).handleMinAge,
		`Mit diesem Angebot erhalten Sie einen einmaligen Rabatt von (\d+)€ auf Ihre monatliche Rechnung\.`:                                                      (*DescriptionParser).handleOneTimeDiscount,
		`Der Mindestbestellwert beträgt (\d+)€\.`:                                                                                                                (*DescriptionParser).handleMinOrderValue,
		`Unsere Techniker kümmern sich um die Installation\.`:                                                                                                    (*DescriptionParser).handleInstallationIncluded,
	}

	// Compile the regex patterns and store them in the map
	compiledPatterns := make(map[*regexp.Regexp]HandlerFunc, len(patters))
	for pattern, HandlerFunc := range patters {
		adjustedPattern := `(?i)` + strings.ReplaceAll(pattern, " ", `\s+`)
		compiledPattern := regexp.MustCompile(adjustedPattern) // Case-insensitive matching
		compiledPatterns[compiledPattern] = HandlerFunc
	}

	return &DescriptionParser{
		emptyPattern: regexp.MustCompile(`^\s*$`),
		patterns:     compiledPatterns,
		logger:       logger,
	}
}

func (d *DescriptionParser) parse(description string) (DescriptionData, error) {
	var data DescriptionData
	for pattern, handleFunc := range d.patterns {
		matches := pattern.FindAllStringSubmatchIndex(description, -1)
		if len(matches) == 0 {
			d.logger.Debug("No match found", "pattern", pattern)
			continue
		}

		if len(matches) > 1 {
			return DescriptionData{}, fmt.Errorf("%w: %s", ErrMultipleMatches, pattern)
		}
		// now, there is exactly one match
		match := matches[0]

		if len(match) < 2 {
			return DescriptionData{}, fmt.Errorf("%w: %s", ErrInvalidMatchFormat, pattern)
		}
		d.logger.Debug("Match found", "pattern", pattern, "match", match, "description", description[match[0]:match[1]])

		// exact one match of the pattern exists in the description
		// use the parse function to extract the data
		err := handleFunc(d, match, description, &data)
		if err != nil {
			return DescriptionData{}, fmt.Errorf("error parsing description: %w", err)
		}

		// remove the matched part from the description
		from := matches[0][0]
		to := matches[0][1]
		description = description[:from] + description[to:]
	}
	d.logger.Debug("Parsed data", "data", data)
	d.logger.Debug("Remaining description", "description", description)

	// Check if the description is empty after parsing
	if !d.emptyPattern.MatchString(description) {
		return DescriptionData{}, fmt.Errorf("could not fully match description data. unmatched: %s", description)
	}

	return data, nil
}

func (d *DescriptionParser) noOp(match []int, description string, target *DescriptionData) error {
	return nil
}

func parseUInt(from int, to int, description string, fieldName string) (int32, error) {
	value, err := strconv.ParseInt(description[from:to], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing %s: %w", fieldName, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("error parsing %s: value cannot be negative", fieldName)
	}
	// valid, as parseUint's bitwidth is set to 32
	return int32(value), nil
}

// parseAndSetUInt parses a substring of the description and sets the value to the target pointer.
// `from` and `to` are the start and end indices of the substring and must be within the bounds of the description.
// `description` is the full description string.
// `target` is a pointer to the target variable where the parsed value will be stored.
// `fieldName` is the name of the field being parsed, used for logging.
func parseAndSetOptUInt(from int, to int, description string, target **int32, fieldName string) error {
	value, err := parseUInt(from, to, description, fieldName)
	if err != nil {
		return err
	}

	*target = &value

	return nil
}

func parseAndSetUnitValue(vFrom int, vTo int, uFrom, uTo int, description string, target **UnitValue, fieldName string) error {
	value, err := strconv.ParseInt(description[vFrom:vTo], 10, 32)
	if err != nil {
		return fmt.Errorf("error parsing %s: %w", fieldName, err)
	}
	if value < 0 {
		return fmt.Errorf("error parsing %s: value cannot be negative", fieldName)
	}
	// valid, as parseUint's bitwidth is set to 32
	*target = &UnitValue{
		Value: int32(value),
		Unit:  description[uFrom:uTo],
	}

	return nil
}

func (d *DescriptionParser) parseSingleUIntPatter(match []int, description string, target **int32, fieldName string) error {
	if len(match) != 4 {
		return fmt.Errorf("invalid match format: expected 4 groups, got %d", len(match))
	}
	err := parseAndSetOptUInt(match[2], match[3], description, target, fieldName)
	if err != nil {
		return err
	}

	d.logger.Debug("Parsed value", fieldName, **target)
	return nil
}

func (d *DescriptionParser) handlePriceTypeAndSpeed(match []int, description string, target *DescriptionData) error {
	if len(match) != 10 {
		return fmt.Errorf("invalid match format: expected 10 groups, got %d", len(match))
	}

	err := parseAndSetOptUInt(match[2], match[3], description, &target.Price, "price")
	if err != nil {
		return err
	}

	ct, err := m.NewConnectionTypeFromValue(strings.ToUpper(description[match[4]:match[5]]))
	if err != nil {
		return fmt.Errorf("error creating connection type: %w", err)
	}
	target.ConnectionType = ct

	err = parseAndSetUnitValue(match[6], match[7], match[8], match[9], description, &target.Speed, "speed")
	if err != nil {
		return err
	}

	d.logger.Debug("Parsed price, connectiontype and speed", "price", target.Price, "connectiontype", target.ConnectionType, "speed", target.Speed.Value, "unit", target.Speed.Unit)
	return nil
}

func (d *DescriptionParser) handleMinimalContractDuration(match []int, description string, target *DescriptionData) error {
	if len(match) != 6 {
		return fmt.Errorf("invalid match format: expected 6 groups, got %d", len(match))
	}

	err := parseAndSetUnitValue(match[2], match[3], match[4], match[5], description, &target.MinimalContractDuration, "minimal contract duration")
	if err != nil {
		return err
	}

	d.logger.Debug("Parsed contract duration", "duration", target.MinimalContractDuration.Value, "unit", target.MinimalContractDuration.Unit)
	return nil
}

func (d *DescriptionParser) handlePercentageDiscount(match []int, description string, target *DescriptionData) error {
	if len(match) != 10 {
		return fmt.Errorf("invalid match format: expected 6 groups, got %d", len(match))
	}

	percent, err := parseUInt(match[2], match[3], description, "percentage discount")
	if err != nil {
		return err
	}
	duration, err := parseUInt(match[4], match[5], description, "percentage discount duration")
	if err != nil {
		return fmt.Errorf("error parsing percentage discount duration: %w", err)
	}

	target.PercentageDiscount = &m.PercentageDiscount{
		Percentage:       percent,
		DurationInMonths: &duration,
	}

	if match[6] != -1 && match[7] != -1 {
		if match[8] == -1 || match[9] == -1 {
			return fmt.Errorf("invalid match format: expected 8 groups, got %d", len(match))
		}

		maxDiscount, err := parseUInt(match[8], match[9], description, "max discount")
		if err != nil {
			return fmt.Errorf("error parsing max discount: %w", err)
		}
		maxDiscount = maxDiscount * units.Eur
		target.PercentageDiscount.MaxDiscountInCent = &maxDiscount
	}

	d.logger.Debug("Parsed discount", "discount", target.PercentageDiscount)
	return nil
}

func (d *DescriptionParser) handleLongRunningPrice(match []int, description string, target *DescriptionData) error {
	if len(match) != 6 {
		return fmt.Errorf("invalid match format: expected 6 groups, got %d", len(match))
	}

	subsequentCostStart, err := parseUInt(match[2], match[3], description, "long running price start")
	if err != nil {
		return err
	}
	subsequentCost, err := parseUInt(match[4], match[5], description, "long running price")
	if err != nil {
		return err
	}
	target.SubsequentCost = &m.SubsequentCost{
		StartMonth:        subsequentCostStart,
		MonthlyCostInCent: subsequentCost * units.Eur,
	}

	d.logger.Debug("Parsed long running price", "subsequentCost", *target.SubsequentCost)
	return nil
}

func (d *DescriptionParser) handleUnthrottledCapacity(match []int, description string, target *DescriptionData) error {
	if len(match) != 6 {
		return fmt.Errorf("invalid match format: expected 6 groups, got %d", len(match))
	}

	err := parseAndSetUnitValue(match[2], match[3], match[4], match[5], description, &target.UnthrottledCapacity, "unthrottled capacity")
	if err != nil {
		return err
	}
	target.UnthrottledCapacity.Unit = description[match[4]:match[5]]

	d.logger.Debug("Parsed unthrottled capacity", "unthrottledCapacity", target.UnthrottledCapacity.Value, "unit", target.UnthrottledCapacity.Unit)
	return nil
}

func (d *DescriptionParser) handleInstallationIncluded(match []int, description string, target *DescriptionData) error {
	if len(match) != 2 {
		return fmt.Errorf("invalid match format: expected 2 groups, got %d", len(match))
	}
	target.InstallationIncluded = true

	d.logger.Debug("Parsed installation included", "installationIncluded", target.InstallationIncluded)
	return nil
}

func (d *DescriptionParser) handleMaxAge(match []int, description string, target *DescriptionData) error {
	return d.parseSingleUIntPatter(match, description, &target.MaxAge, "max age")
}

func (d *DescriptionParser) handleMinAge(match []int, description string, target *DescriptionData) error {
	return d.parseSingleUIntPatter(match, description, &target.MinAge, "min age")
}

func (d *DescriptionParser) handleAdditionalTVChannels(match []int, description string, target *DescriptionData) error {
	if len(match) != 4 {
		return fmt.Errorf("invalid match format: expected 4 groups, got %d", len(match))
	}

	target.IncludedTVSender = description[match[2]:match[3]]

	d.logger.Debug("Parsed additional TV channels", "additionalTVChannels", target.IncludedTVSender)
	return nil
}

func (d *DescriptionParser) handleOneTimeDiscount(match []int, description string, target *DescriptionData) error {
	if len(match) != 4 {
		return fmt.Errorf("invalid match format: expected 4 groups, got %d", len(match))
	}
	value, err := parseUInt(match[2], match[3], description, "one time discount")
	if err != nil {
		return fmt.Errorf("error parsing one time discount: %w", err)
	}
	target.AbsoluteDiscount = &m.AbsoluteDiscount{
		ValueInCent: value * units.Eur,
	}
	d.logger.Debug("Parsed one time discount", "oneTimeDiscount", target.AbsoluteDiscount)
	return nil
}

func (d *DescriptionParser) handleMinOrderValue(match []int, description string, target *DescriptionData) error {
	return d.parseSingleUIntPatter(match, description, &target.MinOrderValue, "min order value")
}
