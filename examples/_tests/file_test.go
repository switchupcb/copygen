package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// checkwd checks the working directory for the test.
func checkwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting the current working directory.\n%v", err)
	}

	if filepath.Base(cwd) != "copygen" {
		// go test uses the package directory as the current working directory.
		if err = os.Chdir(filepath.Join(cwd, "../../")); err != nil {
			t.Fatalf("error changing the current working directory.\n%v", err)
		}
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
