package tests

import (
	"testing"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgoform"
)

// testCommandComparisons represents parameters used to test application commands comparisons.
type testCommandComparisons struct {
	name     string
	a        disgo.CreateGlobalApplicationCommand
	b        disgo.CreateGlobalApplicationCommand
	expected bool
}

// TestCommandComparisons tests application commands comparisons.
func TestCommandComparisons(t *testing.T) {
	tests := []testCommandComparisons{
		{
			name: "basic",
			a: disgo.CreateGlobalApplicationCommand{
				NameLocalizations:        &map[string]string{},
				Description:              new(string),
				DescriptionLocalizations: &map[string]string{},
				DefaultMemberPermissions: nil,
				Type:                     disgo.Pointer(disgo.FlagApplicationCommandTypeCHAT_INPUT),
				NSFW:                     disgo.Pointer(false),
				Name:                     "",
				Options:                  []*disgo.ApplicationCommandOption{},
				IntegrationTypes:         []disgo.Flag{},
				Contexts:                 []disgo.Flag{},
			},
			b:        disgo.CreateGlobalApplicationCommand{},
			expected: false,
		},
		{
			name: "pointer",
			a: disgo.CreateGlobalApplicationCommand{
				Description: new(string),
			},
			b: disgo.CreateGlobalApplicationCommand{
				Description: disgo.Pointer(""),
			},
			expected: true,
		},
		{
			name: "localization-missing-name-localization-key",
			a: disgo.CreateGlobalApplicationCommand{
				Name: "hello",
				NameLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "hello",
					disgo.FlagLocalesEnglishUK:           "mate",
					disgo.FlagLocalesChineseChina:        "你好",
					disgo.FlagLocalesHindi:               "नमस्ते",
					disgo.FlagLocalesSpanish:             "hola",
					disgo.FlagLocalesFrench:              "bonjour",
					disgo.FlagLocalesRussian:             "привет",
					disgo.FlagLocalesPortugueseBrazilian: "olá",
				},
				Description: disgo.Pointer("Say hello."),
				DescriptionLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "Say hello.",
					disgo.FlagLocalesEnglishUK:           "Say hello.",
					disgo.FlagLocalesChineseChina:        "问好。",
					disgo.FlagLocalesHindi:               "नमस्ते बोलो।",
					disgo.FlagLocalesSpanish:             "Di hola.",
					disgo.FlagLocalesFrench:              "Dis bonjour.",
					disgo.FlagLocalesRussian:             "Скажи привет.",
					disgo.FlagLocalesPortugueseBrazilian: "Diga olá.",
				},
				Options: nil,
			},
			b: disgo.CreateGlobalApplicationCommand{
				Name: "hello",
				NameLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUK:           "mate",
					disgo.FlagLocalesChineseChina:        "你好",
					disgo.FlagLocalesHindi:               "नमस्ते",
					disgo.FlagLocalesSpanish:             "hola",
					disgo.FlagLocalesFrench:              "bonjour",
					disgo.FlagLocalesRussian:             "привет",
					disgo.FlagLocalesPortugueseBrazilian: "olá",
				},
				Description: disgo.Pointer("Say hello."),
				DescriptionLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "Say hello.",
					disgo.FlagLocalesEnglishUK:           "Say hello.",
					disgo.FlagLocalesChineseChina:        "问好。",
					disgo.FlagLocalesHindi:               "नमस्ते बोलो।",
					disgo.FlagLocalesSpanish:             "Di hola.",
					disgo.FlagLocalesFrench:              "Dis bonjour.",
					disgo.FlagLocalesRussian:             "Скажи привет.",
					disgo.FlagLocalesPortugueseBrazilian: "Diga olá.",
				},
				Options: nil,
			},
			expected: false,
		},
		{
			name: "localization-diff-name-localization-value",
			a: disgo.CreateGlobalApplicationCommand{
				Name: "hello",
				NameLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "hello",
					disgo.FlagLocalesEnglishUK:           "mate",
					disgo.FlagLocalesChineseChina:        "你好",
					disgo.FlagLocalesHindi:               "नमस्ते",
					disgo.FlagLocalesSpanish:             "hola",
					disgo.FlagLocalesFrench:              "bonjour",
					disgo.FlagLocalesRussian:             "привет",
					disgo.FlagLocalesPortugueseBrazilian: "olá",
				},
				Description: disgo.Pointer("Say hello."),
				DescriptionLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "Say hello.",
					disgo.FlagLocalesEnglishUK:           "Say hello.",
					disgo.FlagLocalesChineseChina:        "问好。",
					disgo.FlagLocalesHindi:               "नमस्ते बोलो।",
					disgo.FlagLocalesSpanish:             "Di hola.",
					disgo.FlagLocalesFrench:              "Dis bonjour.",
					disgo.FlagLocalesRussian:             "Скажи привет.",
					disgo.FlagLocalesPortugueseBrazilian: "Diga olá.",
				},
				Options: nil,
			},
			b: disgo.CreateGlobalApplicationCommand{
				Name: "hello",
				NameLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "abc",
					disgo.FlagLocalesEnglishUK:           "mate",
					disgo.FlagLocalesChineseChina:        "你好",
					disgo.FlagLocalesHindi:               "नमस्ते",
					disgo.FlagLocalesSpanish:             "hola",
					disgo.FlagLocalesFrench:              "bonjour",
					disgo.FlagLocalesRussian:             "привет",
					disgo.FlagLocalesPortugueseBrazilian: "olá",
				},
				Description: disgo.Pointer("Say hello."),
				DescriptionLocalizations: &map[string]string{
					disgo.FlagLocalesEnglishUS:           "Say hello.",
					disgo.FlagLocalesEnglishUK:           "Say hello.",
					disgo.FlagLocalesChineseChina:        "问好。",
					disgo.FlagLocalesHindi:               "नमस्ते बोलो।",
					disgo.FlagLocalesSpanish:             "Di hola.",
					disgo.FlagLocalesFrench:              "Dis bonjour.",
					disgo.FlagLocalesRussian:             "Скажи привет.",
					disgo.FlagLocalesPortugueseBrazilian: "Diga olá.",
				},
				Options: nil,
			},
			expected: false,
		},
	}

	for _, test := range tests {
		if got := disgoform.Equal(test.a, test.b); got != test.expected {
			t.Errorf("(%v: got %v, wanted %v", test.name, got, test.expected)
		}
	}
}
