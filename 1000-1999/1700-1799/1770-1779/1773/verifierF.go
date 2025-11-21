package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	a int
	b int
}

type testInput struct {
	text string
	tc   testCase
}

func buildBinary(path string) (string, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	tmp, err := os.CreateTemp("", "ref1773F")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, base)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build %s failed: %v\n%s", base, err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

type solution struct {
	draws  int
	scores [][2]int
}

func parseOutput(output string, n int) (solution, error) {
	var res solution
	tokens := strings.Fields(output)
	if len(tokens) != n+1 {
		return res, fmt.Errorf("expected %d tokens, got %d", n+1, len(tokens))
	}
	d, err := strconv.Atoi(tokens[0])
	if err != nil {
		return res, fmt.Errorf("invalid draw count: %v", err)
	}
	res.draws = d
	res.scores = make([][2]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Split(tokens[i+1], ":")
		if len(parts) != 2 {
			return res, fmt.Errorf("invalid match token %q", tokens[i+1])
		}
		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return res, fmt.Errorf("invalid score x in token %q: %v", tokens[i+1], err)
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return res, fmt.Errorf("invalid score y in token %q: %v", tokens[i+1], err)
		}
		if x < 0 || y < 0 {
			return res, fmt.Errorf("negative score in token %q", tokens[i+1])
		}
		res.scores[i] = [2]int{x, y}
	}
	return res, nil
}

func verifySolution(tc testCase, sol solution) error {
	if len(sol.scores) != tc.n {
		return fmt.Errorf("expected %d matches, got %d", tc.n, len(sol.scores))
	}
	sumA, sumB := 0, 0
	actualDraws := 0
	for _, sc := range sol.scores {
		sumA += sc[0]
		sumB += sc[1]
		if sc[0] == sc[1] {
			actualDraws++
		}
	}
	if sumA != tc.a {
		return fmt.Errorf("sum of scored goals %d != %d", sumA, tc.a)
	}
	if sumB != tc.b {
		return fmt.Errorf("sum of conceded goals %d != %d", sumB, tc.b)
	}
	if actualDraws != sol.draws {
		return fmt.Errorf("reported draws %d != actual %d", sol.draws, actualDraws)
	}
	if sol.draws < 0 || sol.draws > tc.n {
		return fmt.Errorf("draw count %d out of range", sol.draws)
	}
	return nil
}

func fixedTests() []testInput {
	return []testInput{
		buildInput(testCase{n: 3, a: 2, b: 4}),
		buildInput(testCase{n: 1, a: 2, b: 2}),
		buildInput(testCase{n: 4, a: 0, b: 7}),
		buildInput(testCase{n: 6, a: 3, b: 1}),
	}
}

func buildInput(tc testCase) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n%d\n", tc.n, tc.a, tc.b))
	return testInput{text: sb.String(), tc: tc}
}

func randomTests() []testInput {
	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(100) + 1
		a := rng.Intn(1001)
		b := rng.Intn(1001)
		tests = append(tests, buildInput(testCase{n: n, a: a, b: b}))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierF.go /path/to/candidate\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildBinary(filepath.Join("1000-1999", "1700-1799", "1770-1779", "1773", "1773F.go"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := randomTests()
	for idx, test := range tests {
		refOut, err := runBinary(refBin, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		refSol, err := parseOutput(refOut, test.tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candidate, test.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		candSol, err := parseOutput(candOut, test.tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if err := verifySolution(test.tc, candSol); err != nil {
			fmt.Fprintf(os.Stderr, "candidate invalid on test %d: %v\ninput:\n%s\n", idx+1, err, test.text)
			os.Exit(1)
		}
		if candSol.draws != refSol.draws {
			fmt.Fprintf(os.Stderr, "candidate non-optimal on test %d: draws %d vs optimal %d\ninput:\n%s\n", idx+1, candSol.draws, refSol.draws, test.text)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
