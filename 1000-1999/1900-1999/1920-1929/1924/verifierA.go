package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "cand*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildOracle() (string, func(), error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1924A.go")
	tmp, err := os.CreateTemp("", "oracle*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), src).CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	if k > 26 {
		k = 26
	}
	m := rng.Intn(20) + n
	letters := "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < m; i++ {
		sb.WriteByte(letters[rng.Intn(k)])
	}
	return fmt.Sprintf("1\n%d %d %d\n%s\n", n, k, m, sb.String())
}

// isSubsequence returns true if a is a subsequence of b.
func isSubsequence(a, b string) bool {
	j := 0
	for i := 0; i < len(b) && j < len(a); i++ {
		if b[i] == a[j] {
			j++
		}
	}
	return j == len(a)
}

// validateNoAnswer checks that the candidate's NO answer is valid:
// length n, only first k letters, and not a subsequence of s.
func validateNoAnswer(ans, s string, n, k int) error {
	if len(ans) != n {
		return fmt.Errorf("answer length %d != expected %d", len(ans), n)
	}
	for _, c := range ans {
		if c < 'a' || c >= rune('a'+k) {
			return fmt.Errorf("answer uses invalid character %c (first %d letters only)", c, k)
		}
	}
	if isSubsequence(ans, s) {
		return fmt.Errorf("answer %q is actually a subsequence of %q", ans, s)
	}
	return nil
}

// parseCase extracts n, k, m, s from a single-test-case input string.
func parseCase(input string) (n, k int, s string, err error) {
	lines := strings.Fields(input)
	// lines[0] = "1" (number of test cases)
	// lines[1] = n, lines[2] = k, lines[3] = m
	// lines[4] = s
	if len(lines) < 5 {
		err = fmt.Errorf("unexpected input format")
		return
	}
	n64, e1 := strconv.ParseInt(lines[1], 10, 64)
	k64, e2 := strconv.ParseInt(lines[2], 10, 64)
	if e1 != nil || e2 != nil {
		err = fmt.Errorf("parse error: %v %v", e1, e2)
		return
	}
	n = int(n64)
	k = int(k64)
	s = lines[4]
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	oracle, ocleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ocleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		expUpper := strings.ToUpper(strings.SplitN(exp, "\n", 2)[0])
		gotUpper := strings.ToUpper(strings.SplitN(got, "\n", 2)[0])

		if expUpper != gotUpper {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}

		if gotUpper == "NO" {
			// Validate the candidate's answer instead of comparing to oracle's specific string.
			gotLines := strings.SplitN(got, "\n", 2)
			if len(gotLines) < 2 || strings.TrimSpace(gotLines[1]) == "" {
				fmt.Fprintf(os.Stderr, "case %d: NO answer missing string\ngot:\n%s\ninput:\n%s", i+1, got, input)
				os.Exit(1)
			}
			n, k, s, perr := parseCase(input)
			if perr != nil {
				fmt.Fprintf(os.Stderr, "case %d: parse error: %v\ninput:\n%s", i+1, perr, input)
				os.Exit(1)
			}
			if verr := validateNoAnswer(strings.TrimSpace(gotLines[1]), s, n, k); verr != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid NO answer: %v\ngot:\n%s\ninput:\n%s", i+1, verr, got, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
