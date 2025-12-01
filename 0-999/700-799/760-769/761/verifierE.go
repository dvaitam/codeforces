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
	refSourceE  = "./761E.go"
	randomCases = 200
)

type edge struct {
	u int
	v int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceE)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicTrees()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomCases; i++ {
		tests = append(tests, randomTree(rng))
	}

	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on test %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if normalize(expect) != normalize(got) {
			fail("test %d mismatch\ninput:\n%s\nexpected:\n%s\ngot:\n%s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "761E-ref-*")
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

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func deterministicTrees() []string {
	var tests []string
	tests = append(tests, buildInput(1, nil))
	tests = append(tests, buildInput(2, []edge{{1, 2}}))
	tests = append(tests, buildInput(4, []edge{{1, 2}, {1, 3}, {1, 4}}))
	tests = append(tests, buildInput(5, []edge{{1, 2}, {2, 3}, {3, 4}, {4, 5}}))
	tests = append(tests, buildInput(6, []edge{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}}))
	full := make([]edge, 0, 29*4)
	for i := 2; i <= 30; i++ {
		full = append(full, edge{1, i})
	}
	tests = append(tests, buildInput(30, full))
	return tests
}

func randomTree(rng *rand.Rand) string {
	n := randRange(rng, 1, 30)
	if n == 1 {
		return buildInput(1, nil)
	}
	edges := make([]edge, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, edge{u, v})
	}
	if rng.Intn(4) == 0 {
		shuffleEdges(rng, edges)
	}
	for i := range edges {
		if rng.Intn(2) == 0 {
			edges[i].u, edges[i].v = edges[i].v, edges[i].u
		}
	}
	return buildInput(n, edges)
}

func buildInput(n int, edges []edge) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func shuffleEdges(rng *rand.Rand, edges []edge) {
	rng.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
}

func randRange(rng *rand.Rand, lo, hi int) int {
	if hi < lo {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Intn(hi-lo+1)
}

func normalize(out string) string {
	lines := strings.Fields(out)
	if len(lines) == 0 {
		return ""
	}
	if strings.ToUpper(lines[0]) == "NO" {
		return "NO"
	}
	// otherwise compare everything verbatim because reference output is canonical
	return strings.Join(strings.Fields(out), " ")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
