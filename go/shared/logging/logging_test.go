package logging

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Helper function to capture standard output
func captureOutput(f func()) string {
	// Redirect stdout to capture log output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		logLevel string
		expected string
		err      error
	}{
		{"DEBUG", "DEBUG", nil},
		{"INFO", "INFO", nil},
		{"WARNING", "WARNING", nil},
		{"ERROR", "ERROR", nil},
		{"FATAL", "ERROR", nil}, // FATAL should map to ERROR
		{"INVALID", "", errors.New("Could not create logger, log level INVALID is not supported")},
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			logger, err := NewLogger(tt.logLevel)

			if err != nil && tt.err == nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if err == nil && tt.err != nil {
				t.Fatalf("expected error, got none")
			}
			if err != nil && tt.err != nil {
				if !strings.Contains(err.Error(), tt.err.Error()) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}
			}
			if logger != nil && tt.expected != "" {
				// Logger should not be nil for valid log levels
				if logger.logger == nil {
					t.Fatal("expected a valid slog.Logger, got nil")
				}
			}
		})
	}
}

func TestDebug(t *testing.T) {

	output := captureOutput(func() {
		logger, err := NewLogger("DEBUG")
		if err != nil {
			t.Fatalf("failed to create logger: %v", err)
		}
		logger.Debug("This is a %s message", "debug")
	})

	if !strings.Contains(output, `This is a debug message`) {
		t.Errorf("expected output to contain 'This is a debug message', got %s", output)
	}

	if !strings.Contains(output, "level=DEBUG") {
		t.Errorf("expected output to contain 'level=DEBUG', got %s", output)
	}
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		logger, err := NewLogger("INFO")
		if err != nil {
			t.Fatalf("failed to create logger: %v", err)
		}
		logger.Info("This is an %s message", "info")
	})

	if !strings.Contains(output, "This is an info message") {
		t.Errorf("expected output to contain 'This is an info message', got %s", output)
	}
}

func TestWarning(t *testing.T) {
	output := captureOutput(func() {
		logger, err := NewLogger("WARNING")
		if err != nil {
			t.Fatalf("failed to create logger: %v", err)
		}
		logger.Warning("This is a %s message", "warning")
	})

	if !strings.Contains(output, "This is a warning message") {
		t.Errorf("expected output to contain 'This is a warning message', got %s", output)
	}
}

func TestError(t *testing.T) {
	output := captureOutput(func() {
		logger, err := NewLogger("ERROR")
		if err != nil {
			t.Fatalf("failed to create logger: %v", err)
		}
		logger.Error("This is an %s message", "error")
	})

	if !strings.Contains(output, "This is an error message") {
		t.Errorf("expected output to contain 'This is an error message', got %s", output)
	}
}

func TestFatal(t *testing.T) {
	if os.Getenv("FATAL_CRASHER") == "1" {
		logger, err := NewLogger("FATAL")
		if err != nil {
			t.Fatalf("Failed to create logger: %v", err)
		}
		logger.Fatal("Fatal message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "FATAL_CRASHER=1")
	output, err := cmd.CombinedOutput()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		// The program has exited with an exit code != 0
		// This is expected, so we don't treat it as a test failure
	} else {
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	expected := "Fatal message"
	if !strings.Contains(string(output), expected) {
		t.Errorf("Expected log output to contain %q, but got %q", expected, string(output))
	}
}
