package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	id    string
	input string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}

	base := currentDir()
	refBin, err := buildReference(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, exp, got)
			os.Exit(1)
		}
		if (i+1)%15 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file path")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref2151A.bin")
	cmd := exec.Command("go", "build", "-o", out, "2151A.go")
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		makeCase("single1", 1, []int{1}),
		makeCase("singleLarge", 5, []int{5}),
		makeCase("increasing", 5, []int{1, 2, 3}),
		makeCase("resetOnce", 5, []int{3, 1, 2, 3, 4, 1}),
		makeCase("resetMultiple", 7, []int{4, 1, 2, 3, 1, 2, 3, 4, 5, 1}),
		makeCase("edgeNoFit", 3, []int{2, 3, 4}),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 70; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("rand-%02d", i+1)))
	}

	// larger n with smaller m to stress boundaries
	for i := 0; i < 10; i++ {
		n := 100000
		m := rng.Intn(200) + 1
		tests = append(tests, makeCase(fmt.Sprintf("large-%02d", i+1), n, randomPattern(rng, m)))
	}
	return tests
}

func makeCase(id string, n int, arr []int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{id: id, input: sb.String()}
}

func randomCase(rng *rand.Rand, id string) testCase {
	n := rng.Intn(100000) + 1
	m := rng.Intn(200) + 1
	pattern := randomPattern(rng, m)
	return makeCase(id, n, pattern)
}

func randomPattern(rng *rand.Rand, m int) []int {
	pattern := make([]int, m)
	current := rng.Intn(50) + 1
	for i := 0; i < m; i++ {
		if rng.Float64() < 0.2 {
			current = 1
		} else if rng.Float64() < 0.5 {
			current++
		} else {
			current = rng.Intn(50) + 1
		}
		pattern[i] = max(1, current)
	}
	return pattern
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
