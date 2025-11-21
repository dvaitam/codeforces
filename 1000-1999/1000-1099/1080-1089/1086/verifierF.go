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
	refSourceF   = "1000-1999/1000-1099/1080-1089/1086/1086F.go"
	randomTrials = 200
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceF)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomTest(rng))
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
	tmp, err := os.CreateTemp("", "1086F-ref-*")
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

func deterministicTests() []string {
	tests := []string{
		buildInput(1, 0, [][2]int64{{0, 0}}),
		buildInput(1, 5, [][2]int64{{10, -10}}),
		buildInput(2, 3, [][2]int64{{0, 0}, {2, 2}}),
		buildInput(3, 10, [][2]int64{{-5, -5}, {5, 5}, {0, 0}}),
		buildInput(5, 20, [][2]int64{{1, 2}, {3, 4}, {5, 6}, {-7, 8}, {9, -10}}),
		buildInput(2, 100000000, [][2]int64{{-100000000, 100000000}, {100000000, -100000000}}),
	}
	return tests
}

func randomTest(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	var t int64
	switch rng.Intn(5) {
	case 0:
		t = int64(rng.Intn(50))
	case 1:
		t = int64(rng.Intn(1000))
	case 2:
		t = int64(rng.Intn(100000))
	default:
		t = int64(rng.Intn(100000000))
	}
	pts := make([][2]int64, n)
	for i := 0; i < n; i++ {
		pts[i][0] = randomCoord(rng)
		pts[i][1] = randomCoord(rng)
	}
	return buildInput(n, t, pts)
}

func randomCoord(rng *rand.Rand) int64 {
	switch rng.Intn(3) {
	case 0:
		return int64(rng.Intn(2001) - 1000)
	case 1:
		return int64(rng.Intn(2000001) - 1000000)
	default:
		return int64(rng.Intn(200000001) - 100000000)
	}
}

func buildInput(n int, t int64, pts [][2]int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, t)
	for _, p := range pts {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
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
