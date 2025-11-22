package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refD.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2095D.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validFloatToken(tok string, min, max float64) (float64, error) {
	// ensure at most 6 digits after decimal point
	if strings.Count(tok, ".") > 1 {
		return 0, fmt.Errorf("invalid format %q", tok)
	}
	if idx := strings.IndexByte(tok, '.'); idx != -1 {
		if len(tok)-idx-1 > 6 {
			return 0, fmt.Errorf("too many digits after decimal in %q", tok)
		}
	}
	val, err := strconv.ParseFloat(tok, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %q: %v", tok, err)
	}
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, fmt.Errorf("non-finite value %q", tok)
	}
	if val < min-1e-9 || val > max+1e-9 {
		return 0, fmt.Errorf("value %q out of range [%g,%g]", tok, min, max)
	}
	return val, nil
}

func checkOutput(out string) error {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 numbers, got %d tokens", len(fields))
	}
	if _, err := validFloatToken(fields[0], -90, 90); err != nil {
		return fmt.Errorf("first number invalid: %v", err)
	}
	if _, err := validFloatToken(fields[1], -180, 180); err != nil {
		return fmt.Errorf("second number invalid: %v", err)
	}
	return nil
}

func main() {
	exitCode := 0
	cleanup := func() {}
	defer func() {
		cleanup()
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		exitCode = 1
		return
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	cleanup = func() { _ = os.Remove(ref) }

	if refOut, err := runBinary(ref); err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		exitCode = 1
		return
	} else if err := checkOutput(refOut); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\n", err)
		exitCode = 1
		return
	}

	out, err := runBinary(bin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	if err := checkOutput(out); err != nil {
		fmt.Fprintf(os.Stderr, "output invalid: %v\n", err)
		exitCode = 1
		return
	}

	fmt.Println("All tests passed")
}
