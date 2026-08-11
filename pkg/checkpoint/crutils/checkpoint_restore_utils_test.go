package crutils

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRuntimeProbeScript(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "runtime")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCRRuntimeSupportsExternalUsernsRestore(t *testing.T) {
	runtimePath := writeRuntimeProbeScript(t, `if [ "$1" = restore ] && [ "$2" = --help ]; then
	echo "      --external-userns      restore via caller userns"
	exit 0
fi
exit 1
`)

	if !CRRuntimeSupportsExternalUsernsRestore(runtimePath) {
		t.Fatal("expected runtime help containing --external-userns to be supported")
	}
}

func TestCRRuntimeSupportsExternalUsernsRestoreMissing(t *testing.T) {
	runtimePath := writeRuntimeProbeScript(t, `if [ "$1" = restore ] && [ "$2" = --help ]; then
	echo "      --tcp-established      allow open tcp connections"
	exit 0
fi
exit 1
`)

	if CRRuntimeSupportsExternalUsernsRestore(runtimePath) {
		t.Fatal("expected runtime help without --external-userns to be unsupported")
	}
}
