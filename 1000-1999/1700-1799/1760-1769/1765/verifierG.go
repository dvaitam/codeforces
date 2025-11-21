package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	refSource  = "1000-1999/1700-1799/1760-1769/1765/1765G.go"
	maxQueries = 789
)

type secretCase struct {
	name   string
	secret string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	if err := runInteractive(refBin, tests); err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\n", err)
		os.Exit(1)
	}

	if err := runInteractive(candidate, tests); err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1765G-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runInteractive(program string, tests []secretCase) error {
	cmd := commandFor(program)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr: %w", err)
	}

	var stderrBuf bytes.Buffer
	go io.Copy(&stderrBuf, stderr)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start program: %w", err)
	}

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)

	if _, err := fmt.Fprintf(writer, "%d\n", len(tests)); err != nil {
		return fmt.Errorf("failed to write test count: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush test count: %w", err)
	}

	for idx, tc := range tests {
		if _, err := fmt.Fprintf(writer, "%d\n", len(tc.secret)); err != nil {
			return fmt.Errorf("failed to send n for test %d: %w", idx+1, err)
		}
		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush n for test %d: %w", idx+1, err)
		}
		if err := handleCase(reader, writer, tc); err != nil {
			stdin.Close()
			cmd.Wait()
			return fmt.Errorf("test %d (%s): %w\nstderr:\n%s", idx+1, tc.name, err, stderrBuf.String())
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush final data: %w", err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("program exited with error: %v\nstderr:\n%s", err, stderrBuf.String())
	}
	return nil
}

func handleCase(reader *bufio.Reader, writer *bufio.Writer, tc secretCase) error {
	secret := tc.secret
	n := len(secret)
	if n == 0 || secret[0] != '0' {
		return fmt.Errorf("invalid secret configuration")
	}
	prefix := prefixFunction(secret)
	antiprefix := antiPrefix(secret)
	queries := 0

	for {
		tok, err := nextToken(reader)
		if err != nil {
			return fmt.Errorf("failed to read query type: %w", err)
		}
		switch tok {
		case "1", "2":
			idxTok, err := nextToken(reader)
			if err != nil {
				return fmt.Errorf("failed to read index for query %q: %w", tok, err)
			}
			idx, err := strconv.Atoi(idxTok)
			if err != nil {
				return judgeError(writer, "invalid index %q", idxTok)
			}
			if idx < 1 || idx > n {
				return judgeError(writer, "index out of bounds: %d", idx)
			}
			queries++
			if queries > maxQueries {
				return judgeError(writer, "query limit exceeded (%d)", queries)
			}
			var ans int
			if tok == "1" {
				ans = prefix[idx-1]
			} else {
				ans = antiprefix[idx-1]
			}
			if _, err := fmt.Fprintf(writer, "%d\n", ans); err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("failed to flush response: %w", err)
			}
		case "0":
			guess, err := nextToken(reader)
			if err != nil {
				return fmt.Errorf("failed to read guess: %w", err)
			}
			if len(guess) != n {
				return judgeError(writer, "guess length mismatch: got %d expected %d", len(guess), n)
			}
			if guess != secret {
				return judgeError(writer, "incorrect guess: %s", guess)
			}
			if _, err := fmt.Fprintln(writer, 1); err != nil {
				return fmt.Errorf("failed to acknowledge correct guess: %w", err)
			}
			if err := writer.Flush(); err != nil {
				return fmt.Errorf("failed to flush acknowledgement: %w", err)
			}
			return nil
		default:
			return judgeError(writer, "unknown query type %q", tok)
		}
	}
}

func judgeError(writer *bufio.Writer, format string, args ...interface{}) error {
	fmt.Fprintln(writer, -1)
	writer.Flush()
	return fmt.Errorf(format, args...)
}

func nextToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF && sb.Len() > 0 {
				return sb.String(), nil
			}
			return "", err
		}
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
			if sb.Len() > 0 {
				return sb.String(), nil
			}
			continue
		}
		sb.WriteByte(b)
	}
}

func prefixFunction(s string) []int {
	n := len(s)
	pi := make([]int, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && s[i] != s[j] {
			j = pi[j-1]
		}
		if s[i] == s[j] {
			j++
		}
		pi[i] = j
	}
	return pi
}

func antiPrefix(s string) []int {
	n := len(s)
	if n == 0 {
		return nil
	}
	q := make([]int, n)
	for i := 1; i < n; i++ {
		j := q[i-1]
		for j > 0 && s[j] == s[i] {
			j = q[j-1]
		}
		if s[j] != s[i] {
			q[i] = j + 1
		} else {
			q[i] = 0
		}
	}
	return q
}

func buildTests() []secretCase {
	tests := []secretCase{
		{name: "single-zero", secret: "0"},
		{name: "simple", secret: "01"},
		{name: "sample1", secret: "011001"},
		{name: "sample2", secret: "00111"},
		{name: "alternating", secret: "01010101"},
		{name: "all-zero", secret: "0000000"},
	}

	rng := rand.New(rand.NewSource(1765))
	lengths := []int{12, 25, 40, 75, 120}
	for _, n := range lengths {
		tests = append(tests, secretCase{
			name:   fmt.Sprintf("random-%d", n),
			secret: randomBinary(rng, n),
		})
	}
	return tests
}

func randomBinary(rng *rand.Rand, n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	b[0] = '0'
	for i := 1; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}
