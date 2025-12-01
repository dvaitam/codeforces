package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSourceC1  = "./1116C1.go"
	randomTrials = 200
	maxN         = 200
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceC1)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := exhaustiveInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomInput(rng))
	}

	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if normalize(expect) != normalize(got) {
			fail("case %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1116C1-ref-*")
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

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func exhaustiveInputs() []string {
	var tests []string
	for n := 1; n <= 7; n++ {
		total := 1 << n
		for mask := 0; mask < total; mask++ {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%d ", n))
			for i := n - 1; i >= 0; i-- {
				if mask&(1<<i) != 0 {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			sb.WriteByte('\n')
			tests = append(tests, sb.String())
		}
	}
	return tests
}

func randomInput(rng *rand.Rand) string {
	n := 1
	switch rng.Intn(5) {
	case 0:
		n = rng.Intn(7) + 1
	case 1:
		n = rng.Intn(20) + 1
	case 2:
		n = rng.Intn(100) + 1
	default:
		n = rng.Intn(maxN-1) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d ", n))
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
