package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	maxQueriesPerCase = 30
	totalNLimit       = 10000
)

type testCase struct {
	name string
	perm []int
	posN int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()

	xs, err := runFirstPhase(binary, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "first phase failed: %v\n", err)
		os.Exit(1)
	}

	if err := runSecondPhase(binary, tests, xs); err != nil {
		fmt.Fprintf(os.Stderr, "second phase failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runFirstPhase(bin string, tests []testCase) ([]int, error) {
	var input strings.Builder
	fmt.Fprintln(&input, "first")
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.perm))
		writeArray(&input, tc.perm)
	}
	output, err := runBinaryOnce(bin, input.String())
	if err != nil {
		return nil, err
	}
	tokens := strings.Fields(output)
	if len(tokens) != len(tests) {
		return nil, fmt.Errorf("expected %d outputs in first phase, got %d (%q)", len(tests), len(tokens), output)
	}
	results := make([]int, len(tests))
	for i, tok := range tokens {
		if tok != "0" && tok != "1" {
			return nil, fmt.Errorf("test %d: output %q is not 0/1", i+1, tok)
		}
		if tok == "1" {
			results[i] = 1
		}
	}
	return results, nil
}

func runSecondPhase(bin string, tests []testCase, xs []int) error {
	if len(xs) != len(tests) {
		return fmt.Errorf("internal error: mismatched lengths")
	}

	cmd := commandFor(bin)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture stdout: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to capture stderr: %v", err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start candidate: %v", err)
	}

	errCh := make(chan string, 1)
	go func() {
		data, _ := io.ReadAll(stderr)
		errCh <- string(data)
	}()

	writer := bufio.NewWriter(stdin)
	reader := bufio.NewReader(stdout)

	fmt.Fprintln(writer, "second")
	fmt.Fprintln(writer, len(tests))
	writer.Flush()

	orderRNG := rand.New(rand.NewSource(21682168))
	order := orderRNG.Perm(len(tests))
	for orderIdx, caseIdx := range order {
		fmt.Fprintf(writer, "%d %d\n", len(tests[caseIdx].perm), xs[caseIdx])
		writer.Flush()
		if err := simulateCase(reader, writer, tests[caseIdx]); err != nil {
			stdin.Close()
			cmd.Process.Kill()
			cmd.Wait()
			stderrOut := <-errCh
			return fmt.Errorf("case #%d (%s) at order %d: %v\nstderr:\n%s", caseIdx+1, tests[caseIdx].name, orderIdx+1, err, stderrOut)
		}
	}
	writer.Flush()
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		stderrOut := <-errCh
		return fmt.Errorf("candidate second phase exited with error: %v\nstderr:\n%s", err, stderrOut)
	}
	<-errCh
	return nil
}

func simulateCase(reader *bufio.Reader, writer *bufio.Writer, tc testCase) error {
	queries := 0
	for {
		token, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return fmt.Errorf("unexpected EOF while awaiting query or answer")
			}
			return fmt.Errorf("failed to read token: %v", err)
		}
		switch token {
		case "?":
			queries++
			if queries > maxQueriesPerCase {
				return fmt.Errorf("too many queries (> %d)", maxQueriesPerCase)
			}
			lTok, err := readToken(reader)
			if err != nil {
				return fmt.Errorf("failed to read l after '?': %v", err)
			}
			rTok, err := readToken(reader)
			if err != nil {
				return fmt.Errorf("failed to read r after '?': %v", err)
			}
			lVal, err := strconv.Atoi(lTok)
			if err != nil {
				return fmt.Errorf("invalid l value %q", lTok)
			}
			rVal, err := strconv.Atoi(rTok)
			if err != nil {
				return fmt.Errorf("invalid r value %q", rTok)
			}
			if lVal < 1 || rVal < 1 || lVal > rVal || rVal > len(tc.perm) {
				return fmt.Errorf("query [%d,%d] is out of bounds", lVal, rVal)
			}
			res := rangeDiff(tc.perm, lVal, rVal)
			fmt.Fprintln(writer, res)
			writer.Flush()
		case "!":
			posTok, err := readToken(reader)
			if err != nil {
				return fmt.Errorf("failed to read reported position: %v", err)
			}
			posVal, err := strconv.Atoi(posTok)
			if err != nil {
				return fmt.Errorf("reported position %q is not an integer", posTok)
			}
			if posVal < 1 || posVal > len(tc.perm) {
				return fmt.Errorf("reported position %d out of range", posVal)
			}
			if posVal != tc.posN {
				return fmt.Errorf("reported position %d, expected %d", posVal, tc.posN)
			}
			return nil
		default:
			return fmt.Errorf("unexpected token %q", token)
		}
	}
}

func rangeDiff(perm []int, l, r int) int {
	l--
	r--
	minV, maxV := perm[l], perm[l]
	for i := l + 1; i <= r; i++ {
		if perm[i] < minV {
			minV = perm[i]
		}
		if perm[i] > maxV {
			maxV = perm[i]
		}
	}
	return maxV - minV
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if !isSpace(ch) {
			sb.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if isSpace(ch) {
			break
		}
		sb.WriteByte(ch)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func runBinaryOnce(bin, input string) (string, error) {
	cmd := commandFor(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2168))
	fixedSizes := []int{2, 3, 4, 5, 6, 8, 10, 15, 20, 30, 50, 80, 120, 200, 400, 800, 1600, 3200, 1500}
	var tests []testCase
	total := 0
	for idx, n := range fixedSizes {
		if total+n > totalNLimit {
			break
		}
		perm := randomPermutation(n, rng)
		tests = append(tests, newTestCase(fmt.Sprintf("fixed_%d", idx+1), perm))
		total += n
	}
	for total+2 <= totalNLimit {
		remaining := totalNLimit - total
		maxAdd := remaining
		if maxAdd > 600 {
			maxAdd = 600
		}
		n := rng.Intn(maxAdd-1) + 2
		if total+n > totalNLimit {
			break
		}
		perm := randomPermutation(n, rng)
		tests = append(tests, newTestCase(fmt.Sprintf("rand_%d", len(tests)+1), perm))
		total += n
	}
	return tests
}

func randomPermutation(n int, rng *rand.Rand) []int {
	base := rng.Perm(n)
	for i := range base {
		base[i]++
	}
	idx1, idxN := -1, -1
	for i, v := range base {
		if v == 1 {
			idx1 = i
		}
		if v == n {
			idxN = i
		}
	}
	if idx1 == -1 || idxN == -1 {
		panic("permutation missing 1 or n")
	}
	if idx1 > idxN {
		base[idx1], base[idxN] = base[idxN], base[idx1]
	}
	return base
}

func newTestCase(name string, perm []int) testCase {
	pos := -1
	target := len(perm)
	for i, v := range perm {
		if v == target {
			pos = i + 1
			break
		}
	}
	if pos == -1 {
		panic("permutation missing value n")
	}
	return testCase{name: name, perm: perm, posN: pos}
}
