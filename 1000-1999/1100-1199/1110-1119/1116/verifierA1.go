package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const tol = 1e-6

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1116A1-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleA1")
	cmd := exec.Command("go", "build", "-o", path, "1116A1.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseAmplitudes(out string) ([]float64, error) {
	fields := strings.Fields(out)
	if len(fields) != 4 {
		return nil, fmt.Errorf("expected 4 amplitudes, got %d", len(fields))
	}
	amps := make([]float64, 4)
	for i, f := range fields {
		var val float64
		_, err := fmt.Sscan(f, &val)
		if err != nil {
			return nil, fmt.Errorf("invalid amplitude %q: %v", f, err)
		}
		amps[i] = val
	}
	return amps, nil
}

func compareAmps(exp, got []float64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("amplitude length mismatch")
	}
	for i := range exp {
		if math.Abs(exp[i]-got[i]) > tol {
			return fmt.Errorf("amplitude %d mismatch: expected %.9f got %.9f", i, exp[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	expOut, err := runBinary(oracle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	expAmps, err := parseAmplitudes(expOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid oracle output: %v\noutput:\n%s\n", err, expOut)
		os.Exit(1)
	}

	gotOut, err := runBinary(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	gotAmps, err := parseAmplitudes(gotOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target output invalid: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	if err := compareAmps(expAmps, gotAmps); err != nil {
		fmt.Fprintf(os.Stderr, "mismatch: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}
