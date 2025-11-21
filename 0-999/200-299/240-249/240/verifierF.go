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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"240F.go",
		filepath.Join("0-999", "200-299", "240-249", "240", "240F.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 240F.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref240F_%d.bin", time.Now().UnixNano()))
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
	tests = append(tests, formatTest(1, 1, "a", [][2]int{{1, 1}}))
	tests = append(tests, formatTest(3, 1, "abc", [][2]int{{1, 3}}))
	tests = append(tests, formatTest(5, 2, "aaaaa", [][2]int{{1, 5}, {2, 4}}))
	tests = append(tests, formatTest(6, 3, "abccba", [][2]int{{1, 6}, {2, 5}, {3, 4}}))
	tests = append(tests, formatTest(7, 2, "abcabca", [][2]int{{2, 6}, {1, 7}}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		n := rng.Intn(60) + 1
		m := rng.Intn(60) + 1
		sb := strings.Builder{}
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		qs := make([][2]int, m)
		for i := 0; i < m; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			qs[i] = [2]int{l, r}
		}
		tests = append(tests, formatTest(n, m, sb.String(), qs))
	}
	return tests
}

func formatTest(n, m int, s string, qs [][2]int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, m))
	b.WriteString(s)
	b.WriteByte('\n')
	for _, q := range qs {
		b.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return b.String()
}
