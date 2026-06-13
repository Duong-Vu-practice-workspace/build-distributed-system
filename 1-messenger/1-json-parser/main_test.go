package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) (string, string) {
	oldOut := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	oldErr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	f()

	wOut.Close()
	wErr.Close()
	os.Stdout = oldOut
	os.Stderr = oldErr

	var outBuf, errBuf bytes.Buffer
	io.Copy(&outBuf, rOut)
	io.Copy(&errBuf, rErr)
	return outBuf.String(), errBuf.String()
}

func TestMainParsing(t *testing.T) {
	valid := `{"src":"node1","dest":"node2","body":{"type":"test_type","msg":"hello"}}`
	invalid := `{"src":"bad","dest":"x", "body":` // truncated
	input := strings.Join([]string{valid, invalid}, "\n") + "\n"

	// create pipe and set os.Stdin to the read end
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// write input in a goroutine then close the writer so main sees EOF
	go func() {
		_, _ = w.Write([]byte(input))
		w.Close()
	}()

	out, errOut := captureOutput(func() {
		main()
	})

	if !strings.Contains(out, "node1|node2|test_type") {
		t.Fatalf("stdout missing expected parsed output.\nstdout: %q", out)
	}
	if !strings.Contains(errOut, "Error parsing JSON") {
		t.Fatalf("stderr missing expected error message. stderr: %q", errOut)
	}
}
