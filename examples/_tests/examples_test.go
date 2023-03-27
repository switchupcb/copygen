package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/switchupcb/copygen/cli"
	"github.com/switchupcb/copygen/cli/config"
	"github.com/switchupcb/copygen/cli/generator"
	"github.com/switchupcb/copygen/cli/matcher"
	"github.com/switchupcb/copygen/cli/parser"
)

type test struct {
	name     string
	ymlpath  string // ymlpath represents the path to an example's .yml file.
	wantpath string // wantpath represents the path to a verified example's output file.
}

var (
	tests = []test{
		{
			name:     "main",
			ymlpath:  "main/setup/setup.yml",
			wantpath: "main/copygen.go",
		},
		{
			name:     "automatch",
			ymlpath:  "automatch/setup/setup.yml",
			wantpath: "automatch/copygen.go",
		},
		{
			name:     "basic",
			ymlpath:  "basic/setup/setup.yml",
			wantpath: "basic/copygen.go",
		},
		/*
			{
				name:     "cast-assert",
				ymlpath:  "cast/assert/setup.yml",
				wantpath: "cast/assert/copygen.go",
			},
			{
				name:     "cast-convert",
				ymlpath:  "cast/convert/setup.yml",
				wantpath: "cast/convert/copygen.go",
			},
			{
				name:     "cast-depth",
				ymlpath:  "cast/depth/setup.yml",
				wantpath: "cast/depth/copygen.go",
			},
			{
				name:     "cast-expression",
				ymlpath:  "cast/expression/setup.yml",
				wantpath: "cast/expression/copygen.go",
			},
			{
				name:     "cast-function",
				ymlpath:  "cast/function/setup.yml",
				wantpath: "cast/function/copygen.go",
			},
			{
				name:     "cast-property",
				ymlpath:  "cast/property/setup.yml",
				wantpath: "cast/property/copygen.go",
			},
		*/
		/*
			{
				name:     "deepcopy",
				ymlpath:  "deepcopy/setup/setup.yml",
				wantpath: "deepcopy/copygen.go",
			},
		*/
		{
			name:     "error",
			ymlpath:  "error/setup/setup.yml",
			wantpath: "error/copygen.go",
		},
		{
			name:     "map",
			ymlpath:  "map/setup/setup.yml",
			wantpath: "map/copygen.go",
		},
		{
			name:     "tag",
			ymlpath:  "tag/setup/setup.yml",
			wantpath: "tag/copygen.go",
		},
		{
			name:     "alias",
			ymlpath:  "_tests/alias/setup/setup.yml",
			wantpath: "_tests/alias/copygen.go",
		},
		{
			name:     "automap",
			ymlpath:  "_tests/automap/setup/setup.yml",
			wantpath: "_tests/automap/copygen.go",
		},
		{
			name:     "cyclic",
			ymlpath:  "_tests/cyclic/setup/setup.yml",
			wantpath: "_tests/cyclic/copygen.go",
		},
		{
			name:     "duplicate",
			ymlpath:  "_tests/duplicate/setup/setup.yml",
			wantpath: "_tests/duplicate/copygen.go",
		},
		{
			name:     "import",
			ymlpath:  "_tests/import/setup/setup.yml",
			wantpath: "_tests/import/copygen.go",
		},
		{
			name:     "multi",
			ymlpath:  "_tests/multi/setup/setup.yml",
			wantpath: "_tests/multi/copygen.go",
		},
		{
			name:     "same",
			ymlpath:  "_tests/same/setup/setup.yml",
			wantpath: "_tests/same/setup/copygen.go",
		},
	}
)

// TestExamples tests calls cli.Run() in a similar manner to calling the CLI,
// checking for a valid output.
func TestExamples(t *testing.T) {
	checkwd(t)
	for _, test := range tests {
		testExample(t, test)
	}
}

// testExample tests an example using .go, .tmpl, and programmatic methods.
func testExample(t *testing.T, test test) {
	valid, err := ioutil.ReadFile(test.wantpath)
	if err != nil {
		t.Fatalf("error reading file in test %q.\n%v", test.name, err)
	}

	// test the .go method using CLI Run().
	env := cli.Environment{
		YMLPath: test.ymlpath,
		Output:  false,
		Write:   false,
	}

	goCode, err := env.Run()
	if err != nil {
		t.Fatalf("Run(%q) error: %v", test.name, err)
	}

	if !bytes.Equal(normalizeLineBreaks([]byte(goCode)), normalizeLineBreaks(valid)) {
		fmt.Println(goCode)
		t.Fatalf("Run(%v) output not equivalent to %v", test.name, test.wantpath)
	}

	fmt.Println("PASSED:", test.name)

	// skip the custom generator error example for the .tmpl method.
	if test.name == "error" {
		return
	}

	// test the .tmpl method using copygen programmatically.
	tmplcode, err := templateRun(env)
	if err != nil {
		t.Fatalf("Run(%q [tmpl]) error: %v", test.name, err)
	}

	if !bytes.Equal(normalizeLineBreaks([]byte(tmplcode)), normalizeLineBreaks(valid)) {
		fmt.Println(tmplcode)
		t.Fatalf("Run(%v [tmpl]) output not equivalent to %v", test.name, test.wantpath)
	}

	fmt.Println("PASSED:", test.name, "(tmpl)")
}

// templateRun runs copygen programmatically and generates code using a template.
func templateRun(env cli.Environment) (string, error) {
	gen, err := config.LoadYML(env.YMLPath)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if err = parser.Parse(gen); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if err = matcher.Match(gen); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	gen.Tempath = "tmpl/template/generate.tmpl"
	code, err := generator.Generate(gen, env.Output, env.Write)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return code, nil
}
