package tests

import (
	"reflect"
	"testing"

	"github.com/switchupcb/copygen/cli"
	"github.com/switchupcb/copygen/cli/config"
	"github.com/switchupcb/copygen/cli/parser"
)

// TestGeneratorOptions tests whether the Generator Options are parsed from the setup file correctly.
func TestGeneratorOptions(t *testing.T) {
	checkwd(t)

	env := cli.Environment{
		YMLPath: "_tests/option/setup/setup.yml",
		Output:  false,
		Write:   false,
	}

	gen, err := config.LoadYML(env.YMLPath)
	if err != nil {
		t.Fatalf("Options(%q) error: %v", "Generator", err)
	}

	want := "The possibilities are endless."
	if v, ok := gen.Options.Custom["option"]; ok {
		if vs, ok := v.(string); ok {
			if vs != want {
				t.Fatalf("Options(%q) got %q, want %q", "Generator", vs, want)
			}

			return
		}

		t.Fatalf("Options(%q) does not contain a custom option with a string value.", "Generator")
	}

	t.Fatalf("Options(%q) does not contain a custom option.", "Generator")
}

// TestCustomFunctionOptions tests whether custom Function Options are parsed from the setup file correctly.
func TestCustomFunctionOptions(t *testing.T) {
	checkwd(t)

	env := cli.Environment{
		YMLPath: "_tests/option/setup/setup.yml",
		Output:  false,
		Write:   false,
	}

	gen, err := config.LoadYML(env.YMLPath)
	if err != nil {
		t.Fatalf("Options(%q) error: %v", "Function", err)
	}

	if err = parser.Parse(gen); err != nil {
		t.Fatalf("Options(%q) error: %v", "Function", err)
	}

	wanted := []map[string][]string{
		make(map[string][]string), // A
		make(map[string][]string), // B
		{ // C
			"custom": []string{"comment"},
		},
		{ // D
			"type": []string{"basic"},
		},
		{ // E
			"type": []string{"basic"},
		},
		{ // G
			"type": []string{"basic"},
		},
		{ // F
			"type": []string{"alias"},
		},
		make(map[string][]string), // H
		make(map[string][]string), // Z
	}

	for i, function := range gen.Functions {
		if !reflect.DeepEqual(function.Options.Custom, wanted[i]) {
			t.Fatalf("Options(%q) got %q, want %q", "Function", function.Options.Custom, wanted[i])
		}
	}
}
