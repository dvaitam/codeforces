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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	if candidate == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/candidate")
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
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected: %s\nGot: %s\n", tc.id, tc.input, exp, got)
			os.Exit(1)
		}
		if (i+1)%20 == 0 {
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
	out := filepath.Join(dir, "ref1718A2.bin")
	cmd := exec.Command("go", "build", "-o", out, "1718A2.go")
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
		makeCase("all-zero", [][]int{{0}}),
		makeCase("single", [][]int{{5}}),
		makeCase("two-equal", [][]int{{7, 7}}),
		makeCase("two-different", [][]int{{1, 2}}),
		makeCase("simple-three", [][]int{{1, 2, 3}}),
	)
	// multi-test deterministic
	tests = append(tests, testCase{
		id: "multi-small",
		input: formatTests([][]int{
			{1, 1, 1, 1},
			{0, 0, 0},
			{5},
		}),
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tcCount := rng.Intn(5) + 1
		arrs := make([][]int, tcCount)
		totalN := 0
		for j := 0; j < tcCount; j++ {
			n := rng.Intn(15) + 1
			arrs[j] = make([]int, n)
			for k := 0; k < n; k++ {
				arrs[j][k] = rng.Intn(1 << uint(rng.Intn(10)+1))
			}
			totalN += n
		}
		tests = append(tests, testCase{
			id:    fmt.Sprintf("rand-%02d", i+1),
			input: formatTests(arrs),
		})
	}

	// some larger tests focusing on prefix XOR behavior
	for i := 0; i < 10; i++ {
		tests = append(tests, testCase{
			id: fmt.Sprintf("large-%02d", i+1),
			input: formatTests([][]int{
				makeLargeArray(500, rng),
			}),
		})
	}
	return tests
}

func formatTests(arrs [][]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arrs))
	for _, arr := range arrs {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func makeCase(id string, arrs [][]int) testCase {
	return testCase{id: id, input: formatTests(arrs)}
}

func makeLargeArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		// bias to build interesting XOR prefixes
		if rng.Float64() < 0.5 {
			cur ^= rng.Intn(1 << 20)
		} else {
			cur = rng.Intn(1 << 20)
		}
		arr[i] = cur
	}
	return arr
}
