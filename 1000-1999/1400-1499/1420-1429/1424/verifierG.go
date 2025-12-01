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
	// refSourceG points to the local reference solution to avoid GOPATH resolution.
	refSourceG   = "1424G.go"
	randomTrials = 150
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceG)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomCase(rng))
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
	tmp, err := os.CreateTemp("", "1424G-ref-*")
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

func deterministicCases() []string {
	return []string{
		buildInput([][2]int{{1, 2}}),
		buildInput([][2]int{{1, 5}, {2, 6}, {3, 7}}),
		buildInput([][2]int{{10, 20}, {15, 25}, {5, 30}, {1, 40}}),
		buildInput([][2]int{{100, 101}, {101, 102}, {102, 103}}),
		buildInput([][2]int{{1, 3}, {2, 4}, {3, 5}, {4, 6}, {5, 7}}),
	}
}

func randomCase(rng *rand.Rand) string {
	n := randomN(rng)
	intervals := make([][2]int, n)
	for i := 0; i < n; i++ {
		b := randomYear(rng)
		d := b + rng.Intn(1000) + 1
		if rng.Intn(5) == 0 {
			d = b + 1
		}
		if b > d {
			b, d = d, b
		}
		if b == d {
			d++
		}
		intervals[i] = [2]int{b, d}
	}
	return buildInput(intervals)
}

func randomN(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 1
	case 1:
		return rng.Intn(5) + 2
	case 2:
		return rng.Intn(50) + 10
	case 3:
		return rng.Intn(1000) + 100
	default:
		return rng.Intn(100000) + 1000
	}
}

func randomYear(rng *rand.Rand) int {
	switch rng.Intn(4) {
	case 0:
		return rng.Intn(20) + 1
	case 1:
		return rng.Intn(1000) + 1
	case 2:
		return rng.Intn(100000) + 1
	default:
		return rng.Intn(1000000000-1) + 1
	}
}

func buildInput(intervals [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(intervals)))
	for _, iv := range intervals {
		sb.WriteString(fmt.Sprintf("%d %d\n", iv[0], iv[1]))
	}
	return sb.String()
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
