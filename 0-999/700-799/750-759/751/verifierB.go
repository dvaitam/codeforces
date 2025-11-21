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

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"751B.go",
		filepath.Join("0-999", "700-799", "750-759", "751", "751B.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 751B.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref751B_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []testCase {
	tests := []testCase{
		newTest(2, 1, []int{0, 1}),
		newTest(3, 2, []int{-2, 0, 2}),
		newTest(4, 5, []int{10, 15, 0, -5}),
		newTest(5, 1000000000, []int{0, 5, 10, 15, 20}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 160 {
		var n int
		switch {
		case len(tests) < 40:
			n = rng.Intn(8) + 2
		case len(tests) < 80:
			n = rng.Intn(40) + 10
		case len(tests) < 120:
			n = rng.Intn(400) + 100
		default:
			n = rng.Intn(2000) + 500
		}
		if n > 100000 {
			n = 100000
		}
		d := rng.Intn(1_000_000_000) + 1
		coords := make([]int, n)
		used := make(map[int]struct{}, n)
		for i := 0; i < n; i++ {
			val := rng.Intn(2_000_000_001) - 1_000_000_000
			for {
				if _, ok := used[val]; !ok {
					break
				}
				val++
				if val > 1_000_000_000 {
					val = -1_000_000_000
				}
			}
			used[val] = struct{}{}
			coords[i] = val
		}
		tests = append(tests, newTest(n, d, coords))
	}
	return tests
}

func newTest(n, d int, coords []int) testCase {
	if len(coords) != n {
		panic("coordinate count mismatch")
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, d))
	for i, v := range coords {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}
