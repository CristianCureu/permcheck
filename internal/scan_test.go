package internal

import (
	"os"
	"testing"
)

func TestRunScan_FixInsecureFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "permcheck_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create insecure file
	filePath := tmpDir + "/test_insecure.txt"
	err = os.WriteFile(filePath, []byte("test"), 0666) // world writable
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	cfg := &Config{
		FixMode:      true,
		InsecureOnly: false,
		NumWorkers:   1,
		IsTerminal:   false,
		InsecureTag:  DefaultInsecureTag,
		SecureTag:    DefaultSecureTag,
	}

	err = RunScan([]string{tmpDir}, cfg)
	if err != nil {
		t.Fatalf("RunScan failed: %v", err)
	}

	// Check permissions were fixed
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	expectedMode := os.FileMode(0644)
	if info.Mode().Perm() != expectedMode {
		t.Errorf("Permissions not fixed: got %04o, want %04o", info.Mode().Perm(), expectedMode)
	}
}
