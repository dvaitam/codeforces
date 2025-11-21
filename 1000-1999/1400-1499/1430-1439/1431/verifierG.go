package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	refSourceG   = "1000-1999/1400-1499/1430-1439/1431/1431G.go"
	randomTrials = 200
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
	tmp, err := os.CreateTemp("", "1431G-ref-*")
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
	cases := []string{
		buildInput([]int{1, 2}, 1),
		buildInput([]int{1, 3, 5, 7}, 2),
		buildInput([]int{10, 20, 30, 40, 50, 60}, 3),
		buildInput([]int{5, 4, 3, 2, 1, 0}, 2),
		buildInput([]int{1, 1000000}, 1),
	}
	return cases
}

func randomCase(rng *rand.Rand) string {
	n := randomN(rng)
	k := rng.Intn(n/2) + 1
	values := randomValues(rng, n)
	return buildInput(values, k)
}

func randomN(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 2
	case 1:
		return rng.Intn(5) + 3
	case 2:
		return rng.Intn(40) + 10
	case 3:
		return rng.Intn(150) + 50
	default:
		return rng.Intn(400-200) + 200
	}
}

func randomValues(rng *rand.Rand, n int) []int {
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(1_000_000) + 1
	}
	sort.Ints(vals)
	return vals
}

func buildInput(values []int, k int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", len(values), k))
	for i, v := range values {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
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
