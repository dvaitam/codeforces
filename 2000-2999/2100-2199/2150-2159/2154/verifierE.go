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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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
		if (i+1)%10 == 0 {
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
	out := filepath.Join(dir, "ref2154E.bin")
	cmd := exec.Command("go", "build", "-o", out, "2154E.go")
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
		makeCase("simple", []int{5}, []int{1, 1, 5, 5, 5}, 1),
		makeCase("k1", []int{3}, []int{10, 1, 2}, 1),
		makeCase("kLarge", []int{6}, []int{1, 1, 2, 3, 5, 8}, 3),
		makeCase("allSame", []int{5}, []int{7, 7, 7, 7, 7}, 2),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// random tests respecting total n <= 2e5
	totalN := 0
	for i := 0; i < 60 && totalN < 200000; i++ {
		n := rng.Intn(50) + 3
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		k := rng.Intn(100) + 1
		arr := randomArray(rng, n)
		tests = append(tests, makeCase(fmt.Sprintf("rand-%02d", i+1), arr, nil, k))
		totalN += n
	}
	// stress: large n, sorted arrays, descending, etc.
	tests = append(tests,
		makeCase("sorted", rangeArray(200000), nil, 2),
		makeCase("desc", rangeArrayDesc(200000), nil, 3),
		makeCase("alternating", alternatingArray(200000), nil, 5),
	)
	return tests
}

func makeCase(id string, arr []int, override []int, k int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", len(arr), k)
	if override != nil {
		arr = override
	}
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{id: id, input: sb.String()}
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	base := rng.Intn(1000) + 1
	for i := 0; i < n; i++ {
		if rng.Float64() < 0.3 {
			base = rng.Intn(1000) + 1
		} else {
			base += rng.Intn(5) - 2
			if base < 1 {
				base = 1
			}
		}
		arr[i] = base
	}
	return arr
}

func rangeArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	return arr
}

func rangeArrayDesc(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = n - i
	}
	return arr
}

func alternatingArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			arr[i] = 1
		} else {
			arr[i] = 1000000000
		}
	}
	return arr
}
