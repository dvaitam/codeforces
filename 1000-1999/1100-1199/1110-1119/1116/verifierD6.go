package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSourceD6  = "./1116D6.go"
	randomTrials = 30
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD6.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceD6)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicArgs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, []string{fmt.Sprintf("%d", rng.Intn(3)+2)})
	}

	for idx, args := range tests {
		expect, err := runProgram(refBin, args)
		if err != nil {
			fail("reference failed on case %d: %v\nargs:%v", idx+1, err, args)
		}
		if err := validateMatrix(expect, args); err != nil {
			fail("reference output invalid on case %d: %v\nargs:%v\noutput:\n%s", idx+1, err, args, expect)
		}
		got, err := runCandidate(candidate, args)
		if err != nil {
			fail("candidate failed on case %d: %v\nargs:%v", idx+1, err, args)
		}
		if err := validateMatrix(got, args); err != nil {
			fail("candidate output invalid on case %d: %v\nargs:%v\noutput:\n%s", idx+1, err, args, got)
		}
		if normalize(expect) != normalize(got) {
			fail("case %d mismatch\nargs:%v\nexpected: %s\ngot: %s", idx+1, args, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1116D6-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin string, args []string) (string, error) {
	cmd := exec.Command(bin, args...)
	return runCommand(cmd)
}

func runCandidate(target string, args []string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmdArgs := append([]string{"run", target}, args...)
		cmd = exec.Command("go", cmdArgs...)
	} else {
		cmd = exec.Command(target, args...)
	}
	return runCommand(cmd)
}

func runCommand(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicArgs() [][]string {
	return [][]string{
		nil,
		{"2"},
		{"3"},
		{"4"},
	}
}

func validateMatrix(out string, args []string) error {
	N := 2
	if len(args) > 0 {
		fmt.Sscanf(args[0], "%d", &N)
	}
	if N < 2 || N > 4 {
		return fmt.Errorf("N=%d out of range", N)
	}
	size := 1 << N
	fields := strings.Fields(out)
	if len(fields) != size*size {
		return fmt.Errorf("expected %d entries, got %d", size*size, len(fields))
	}
	const eps = 1e-6
	for col := 0; col < size; col++ {
		for row := 0; row < size; row++ {
			idx := col*size + row
			val, err := parseFloat(fields[idx])
			if err != nil {
				return fmt.Errorf("invalid float %q: %v", fields[idx], err)
			}
			if row < col-1 && math.Abs(val) > eps {
				return fmt.Errorf("entry (%d,%d) should be zero", row, col)
			}
		}
	}
	return nil
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func normalize(out string) string {
	return strings.Join(strings.Fields(out), " ")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
