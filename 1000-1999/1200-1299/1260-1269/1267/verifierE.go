package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const refSource = "1000-1999/1200-1299/1260-1269/1267/1267E.go"

type testCase struct {
	name  string
	input string
	n     int
	m     int
	votes [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refSol, err := parseSolution(tc, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candSol, err := parseSolution(tc, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if candSol.k > refSol.k {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): candidate cancels %d stations, expected <= %d\n", idx+1, tc.name, candSol.k, refSol.k)
			os.Exit(1)
		}
		if !isValidSolution(tc, candSol) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): cancellation list %v is invalid\n", idx+1, tc.name, candSol.indices)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1267E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

type solution struct {
	k       int
	indices []int
}

func parseSolution(tc testCase, out string) (solution, error) {
	lines := filterNonEmpty(strings.Split(out, "\n"))
	if len(lines) == 0 {
		return solution{}, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(strings.Fields(lines[0])[0])
	if err != nil {
		return solution{}, fmt.Errorf("invalid k: %v", err)
	}
	var indices []int
	if k > 0 {
		var tokens []string
		if len(lines) > 1 {
			tokens = strings.Fields(lines[1])
		}
		if len(tokens) != k {
			return solution{}, fmt.Errorf("expected %d indices, got %d", k, len(tokens))
		}
		indices = make([]int, k)
		for i, tok := range tokens {
			val, err := strconv.Atoi(tok)
			if err != nil {
				return solution{}, fmt.Errorf("invalid index %q", tok)
			}
			val--
			if val < 0 || val >= tc.m {
				return solution{}, fmt.Errorf("index %d out of range", val+1)
			}
			indices[i] = val
		}
	}
	return solution{k: k, indices: indices}, nil
}

func filterNonEmpty(lines []string) []string {
	var res []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			res = append(res, line)
		}
	}
	return res
}

func isValidSolution(tc testCase, sol solution) bool {
	used := make([]bool, tc.m)
	for _, idx := range sol.indices {
		if used[idx] {
			return false
		}
		used[idx] = true
	}
	active := make([]bool, tc.m)
	for i := range active {
		active[i] = true
	}
	for _, idx := range sol.indices {
		active[idx] = false
	}
	sums := make([]int, tc.n)
	for i := 0; i < tc.m; i++ {
		if !active[i] {
			continue
		}
		for j := 0; j < tc.n; j++ {
			sums[j] += tc.votes[i][j]
		}
	}
	opposition := sums[tc.n-1]
	for j := 0; j < tc.n-1; j++ {
		if opposition > sums[j] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("simple", [][]int{
			{5, 4},
			{3, 6},
		}),
		buildCase("tie-sum", [][]int{
			{4, 5, 6},
			{5, 4, 6},
		}),
		buildCase("needs-many", [][]int{
			{10, 0, 0, 20},
			{5, 5, 5, 15},
			{0, 10, 0, 15},
			{2, 2, 2, 10},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(5) + 2
		m := rng.Intn(10) + 1
		matrix := make([][]int, m)
		for r := 0; r < m; r++ {
			matrix[r] = make([]int, n)
			for c := 0; c < n; c++ {
				matrix[r][c] = rng.Intn(1000)
			}
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), matrix))
	}
	return tests
}

func buildCase(name string, matrix [][]int) testCase {
	m := len(matrix)
	n := len(matrix[0])
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(matrix[i][j]))
		}
		sb.WriteByte('\n')
	}
	copyVotes := make([][]int, m)
	for i := range copyVotes {
		copyVotes[i] = append([]int(nil), matrix[i]...)
	}
	return testCase{name: name, input: sb.String(), n: n, m: m, votes: copyVotes}
}
