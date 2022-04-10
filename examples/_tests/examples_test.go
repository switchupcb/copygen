package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/switchupcb/copygen/cli"
)

var (
	tests = []struct {
		name     string
		ymlpath  string // ymlpath represents the path to an example's .yml file.
		wantpath string // wantpath represents the path to a verified example's output file.
	}{
		{
			name:     "main",
			ymlpath:  "examples/main/setup/setup.yml",
			wantpath: "examples/main/copygen.go",
		},
		{
			name:     "automatch",
			ymlpath:  "examples/automatch/setup/setup.yml",
			wantpath: "examples/automatch/copygen.go",
		},

		/*
			{
				name:     "deepcopy",
				ymlpath:  "examples/deepcopy/setup/setup.yml",
				wantpath: "examples/deepcopy/copygen.go",
			},
		*/

		{
			name:     "error",
			ymlpath:  "examples/error/setup/setup.yml",
			wantpath: "examples/error/copygen.go",
		},
		{
			name:     "manual",
			ymlpath:  "examples/manual/setup/setup.yml",
			wantpath: "examples/manual/copygen.go",
		},
		{
			name:     "alias",
			ymlpath:  "examples/_tests/alias/setup/setup.yml",
			wantpath: "examples/_tests/alias/copygen.go",
		},
		{
			name:     "cyclic",
			ymlpath:  "examples/_tests/cyclic/setup/setup.yml",
			wantpath: "examples/_tests/cyclic/copygen.go",
		},

		// .tmpl
		{
			name:     "main (tmpl)",
			ymlpath:  "examples/tmpl/setup/setup.yml",
			wantpath: "examples/tmpl/copygen.go",
		},
	}
)

// TestExamples tests calls cli.Run() in a similar manner to calling the CLI,
// checking for a valid output.
func TestExamples(t *testing.T) {
	// go test uses the package directory as the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting the current working directory.\n%v", err)
	}
	if err = os.Chdir(filepath.Join(cwd, "../../")); err != nil {
		t.Fatalf("error changing the current working directory.\n%v", err)
	}

	for _, test := range tests {
		env := cli.Environment{
			YMLPath: test.ymlpath,
			Output:  false,
			Write:   false,
		}

		code, err := env.Run()
		if err != nil {
			t.Fatalf("Run(%q) error: %v", test.name, err)
		}

		valid, err := ioutil.ReadFile(test.wantpath)
		if err != nil {
			t.Fatalf("error reading file in test %q.\n%v", test.name, err)
		}

		if !bytes.Equal(normalizeLineBreaks([]byte(code)), normalizeLineBreaks(valid)) {
			fmt.Println(code)
			t.Fatalf("Run(%v) output not equivalent to %v", test.name, test.wantpath)
		}
		fmt.Println("Passed:", test.name)
	}
}

// normalizeLineBreaks normalizes line breaks for file comparison.
func normalizeLineBreaks(d []byte) []byte {
	// replace CRLF \r\n with LF \n
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r with LF \n
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}
