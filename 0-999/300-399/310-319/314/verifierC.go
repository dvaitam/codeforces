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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	for idx, input := range tests {
		want, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if normalize(got) != normalize(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"314C.go",
		filepath.Join("0-999", "300-399", "310-319", "314", "314C.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 314C.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref314C_%d.bin", time.Now().UnixNano()))
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
	return strings.TrimSpace(stdout.String()), nil
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func generateTests() []string {
	var tests []string
	tests = append(tests, formatTest([]int{1}))
	tests = append(tests, formatTest([]int{5, 5, 5, 5}))
	tests = append(tests, formatTest([]int{1, 2, 3, 4, 5}))
	tests = append(tests, formatTest([]int{5, 4, 3, 2, 1}))
	tests = append(tests, formatTest([]int{1, 1000000, 2, 999999, 3}))
	tests = append(tests, formatLargeConstant(1000, 7))
	tests = append(tests, formatLargeRange(2000))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 150 {
		var n int
		if len(tests)%15 == 0 {
			n = 10000 + rng.Intn(5001) // bigger cases
		} else {
			n = rng.Intn(80) + 1
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			val := rng.Intn(1000000) + 1
			// occasional repeated blocks
			if rng.Intn(5) == 0 && i > 0 {
				val = a[i-1]
			}
			a[i] = val
		}
		tests = append(tests, formatTest(a))
	}
	return tests
}

func formatTest(arr []int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return b.String()
}

func formatLargeConstant(n, value int) string {
	a := make([]int, n)
	for i := range a {
		a[i] = value
	}
	return formatTest(a)
}

func formatLargeRange(n int) string {
	a := make([]int, n)
	cur := 1
	for i := range a {
		a[i] = cur
		cur++
		if cur > 1000000 {
			cur = 1
		}
	}
	return formatTest(a)
}
