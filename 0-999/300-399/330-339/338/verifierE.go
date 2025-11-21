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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		"338E.go",
		filepath.Join("0-999", "300-399", "330-339", "338", "338E.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 338E.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref338E_%d.bin", time.Now().UnixNano()))
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
	tests = append(tests, formatTest(1, 1, 0, []int64{3}, []int64{3}))
	tests = append(tests, formatTest(2, 1, 5, []int64{10}, []int64{15, 20}))
	tests = append(tests, formatTest(5, 3, 2, []int64{1, 3, 5}, []int64{1, 2, 3, 4, 5}))
	tests = append(tests, formatTest(3, 3, 10, []int64{4, 7, 11}, []int64{5, 8, 12}))
	tests = append(tests, formatTest(4, 2, 0, []int64{1000000000, 1}, []int64{1000000000, 2, 1000000000, 3}))
	tests = append(tests, formatIncreasing(2000, 500))
	tests = append(tests, formatConstant(1500, 800))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 160 {
		var n int
		if len(tests)%20 == 0 {
			n = 50000 + rng.Intn(20000)
		} else {
			n = rng.Intn(80) + 1
		}
		m := rng.Intn(n) + 1
		h := rng.Int63n(1_000_000_000) + 1
		b := make([]int64, m)
		a := make([]int64, n)
		for i := 0; i < m; i++ {
			b[i] = rng.Int63n(1_000_000_000) + 1
		}
		for i := 0; i < n; i++ {
			val := rng.Int63n(1_000_000_000) + 1
			if rng.Intn(6) == 0 && i > 0 {
				diff := int64(rng.Intn(2000)) - 1000
				val = a[i-1] + diff
				if val < 1 {
					val = 1
				}
				if val > 1_000_000_000 {
					val = 1_000_000_000
				}
			}
			a[i] = val
		}
		tests = append(tests, formatTest(n, m, h, b, a))
	}
	return tests
}

func formatTest(n, m int, h int64, b, a []int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, h))
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func formatIncreasing(n, m int) string {
	if m > n {
		m = n
	}
	h := int64(5)
	b := make([]int64, m)
	a := make([]int64, n)
	for i := 0; i < m; i++ {
		b[i] = int64(i * 3)
	}
	for i := 0; i < n; i++ {
		a[i] = int64(i * 3)
	}
	return formatTest(n, m, h, b, a)
}

func formatConstant(n, m int) string {
	if m > n {
		m = n
	}
	h := int64(0)
	b := make([]int64, m)
	a := make([]int64, n)
	for i := 0; i < m; i++ {
		b[i] = 1
	}
	for i := 0; i < n; i++ {
		a[i] = 1
	}
	return formatTest(n, m, h, b, a)
}
