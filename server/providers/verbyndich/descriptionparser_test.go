package verbyndich

import (
	"log/slog"
	"testing"

	m "github.com/rotmanjanez/check24-gendev-7/pkg/models"
)

func TestDescriptionParser_WhitespaceVariations(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)

	testCases := []struct {
		name        string
		description string
	}{
		{
			name:        "NormalWhitespace",
			description: "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s.",
		},
		{
			name:        "ExtraSpaces",
			description: "Für  nur  29€  im  Monat  erhalten  Sie  eine  DSL-Verbindung  mit  einer  Geschwindigkeit  von  100  Mbit/s.",
		},
		{
			name:        "LeadingTrailingSpaces",
			description: "  Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s.  ",
		},
		{
			name:        "TabsInsteadOfSpaces",
			description: "Für\tnur\t29€\tim\tMonat\terhalten\tSie\teine\tDSL-Verbindung\tmit\teiner\tGeschwindigkeit\tvon\t100\tMbit/s.",
		},
		{
			name:        "MixedWhitespaceAndNewlines",
			description: "Für nur 29€ im Monat\n erhalten Sie eine DSL-Verbindung \t mit einer Geschwindigkeit von 100 Mbit/s.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.parse(tc.description)
			if err != nil {
				t.Errorf("expected success, got error: %v", err)
			}
		})
	}
}

func TestDescriptionParser_EdgeCaseValues(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)

	testCases := []struct {
		name        string
		description string
		shouldWork  bool
	}{
		{
			name:        "ZeroPrice",
			description: "Für nur 0€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s.",
			shouldWork:  true,
		},
		{
			name:        "ZeroSpeed",
			description: "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 0 Mbit/s.",
			shouldWork:  true,
		},
		{
			name:        "NegativePrice",
			description: "Für nur -29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s.",
			shouldWork:  false,
		},
		{
			name:        "VeryHighSpeed",
			description: "Für nur 99€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 10000 Mbit/s.",
			shouldWork:  true,
		},
		{
			name:        "GbitSpeed",
			description: "Für nur 79€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 1 Gbit/s.",
			shouldWork:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.parse(tc.description)
			if tc.shouldWork && err != nil {
				t.Errorf("expected success, got error: %v", err)
			}
			if !tc.shouldWork && err == nil {
				t.Errorf("expected error, got success")
			}
		})
	}
}

func TestDescriptionParser_UnknownDescription(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "This is an unknown description pattern that should not match any regex."
	_, err := parser.parse(description)
	if err == nil {
		t.Errorf("expected error for unknown description pattern, got success")
	}
}

func TestDescriptionParser_ConflictingInformation(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Für nur 39€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 200 Mbit/s."
	_, err := parser.parse(description)
	if err == nil {
		t.Errorf("expected error for conflicting price statements, got success")
	}
}

func TestDescriptionParser_ConflictingAgeRestrictions(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "Für nur 29€ im Monat erhalten Sie eine DSL-Verbindung mit einer Geschwindigkeit von 100 Mbit/s. Dieses Angebot ist nur für Personen unter 65 Jahren verfügbar. Dieses Angebot ist nur für Personen über 18 Jahren verfügbar."
	_, err := parser.parse(description)
	if err != nil {
		t.Errorf("expected success for non-conflicting age restrictions, got error: %v", err)
	}
}

func TestDescriptionParser_ComplexValidDescription(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "Für nur 49€ im Monat erhalten Sie eine FIBER-Verbindung mit einer Geschwindigkeit von 500 Mbit/s. Bitte beachten Sie, dass die Mindestvertragslaufzeit 24 Monate beträgt. Mit diesem Angebot erhalten Sie einen Rabatt von 10% auf Ihre monatliche Rechnung bis zum 12. Monat. Der maximale Rabatt beträgt 30€. Ab dem 25. Monat beträgt der monatliche Preis 59€. Ab 100GB pro Monat wird die Geschwindigkeit gedrosselt. Zusätzlich sind folgende Fernsehsender enthalten Premium+. Dieses Angebot ist nur für Personen unter 30 Jahren verfügbar. Mit diesem Angebot erhalten Sie einen einmaligen Rabatt von 50€ auf Ihre monatliche Rechnung. Der Mindestbestellwert beträgt 25€. Unsere Techniker kümmern sich um die Installation."
	_, err := parser.parse(description)
	if err != nil {
		t.Errorf("expected success for complex valid description, got error: %v", err)
	}
}

func TestDescriptionParser_PartiallyValidDescription(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "Für nur 35€ im Monat erhalten Sie eine CABLE-Verbindung mit einer Geschwindigkeit von 250 Mbit/s. Hier ist ein unbekannter Satz der ignoriert werden sollte. Zusätzlich sind folgende Fernsehsender enthalten Sports."
	_, err := parser.parse(description)
	if err == nil {
		t.Errorf("expected error for partially valid description with unknown sentence, got success")
	}
}

func TestDescriptionParser_InvalidConnectionType(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "Für nur 29€ im Monat erhalten Sie eine QUANTUM-Verbindung mit einer Geschwindigkeit von 100 Mbit/s."
	_, err := parser.parse(description)
	if err == nil {
		t.Errorf("expected error for unsupported connection type, got success")
	}
}

func TestDescriptionParser_ExcessiveWhitespace(t *testing.T) {
	logger := slog.Default()
	parser := NewDescriptionParser(logger)
	description := "   Für nur 29€ im Monat   erhalten Sie eine DSL-Verbindung mit einer    Geschwindigkeit von 100 Mbit/s.     Zusätzlich sind folgende Fernsehsender enthalten TestTV.   "
	result, err := parser.parse(description)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if result.Speed == nil || result.Speed.Value != 100 || result.Speed.Unit != "Mbit/s" {
		t.Errorf("expected speed 100 Mbit/s, got %+v", result.Speed)
	}
	if result.ConnectionType != m.DSL {
		t.Errorf("expected connection type DSL, got %v", result.ConnectionType)
	}
	if result.IncludedTVSender != "TestTV" {
		t.Errorf("expected IncludedTVSender 'TestTV', got %v", result.IncludedTVSender)
	}
	if result.Price == nil || *result.Price != 29 {
		t.Errorf("expected price 29, got %v", result.Price)
	}
}

// Add more parser-specific tests as needed
