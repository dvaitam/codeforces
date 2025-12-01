package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource  = "./2052A.go"
	maxN       = 1000
	totalLimit = 1000
)

type testCase struct {
	n   int
	per []int
}

type move struct {
	x int
	y int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	candidate := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		refSeq := parseSequence(refOut)

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}
		candSeq := parseSequence(candOut)

		if err := validateSequence(tc, candSeq); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid answer: %v\n", idx+1, err)
			os.Exit(1)
		}
		if len(candSeq) != len(refSeq) {
			fmt.Fprintf(os.Stderr, "test %d wrong length: expected %d moves, got %d\n", idx+1, len(refSeq), len(candSeq))
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2052A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func parseSequence(out string) []move {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	idx := 0
	for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
		idx++
	}
	if idx >= len(lines) {
		fmt.Fprintln(os.Stderr, "missing number of moves in output")
		os.Exit(1)
	}
	mVal, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse number of moves: %v\n", err)
		os.Exit(1)
	}
	idx++
	ops := make([]move, 0, mVal)
	for idx < len(lines) && len(ops) < mVal {
		line := strings.TrimSpace(lines[idx])
		idx++
		if line == "" {
			continue
		}
		var x, y int
		if _, err := fmt.Sscanf(line, "%d %d", &x, &y); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse move line %q: %v\n", line, err)
			os.Exit(1)
		}
		ops = append(ops, move{x: x, y: y})
	}
	if len(ops) != mVal {
		fmt.Fprintf(os.Stderr, "expected %d moves, got %d\n", mVal, len(ops))
		os.Exit(1)
	}
	return ops
}

func validateSequence(tc testCase, ops []move) error {
	n := tc.n
	current := make([]int, n)
	for i := 0; i < n; i++ {
		current[i] = i + 1
	}
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		pos[current[i]] = i
	}
	used := make(map[[2]int]bool)
	for idx, mv := range ops {
		if mv.x < 1 || mv.x > n || mv.y < 1 || mv.y > n {
			return fmt.Errorf("move %d uses invalid cars %d %d", idx+1, mv.x, mv.y)
		}
		if mv.x == mv.y {
			return fmt.Errorf("move %d: x equals y (%d)", idx+1, mv.x)
		}
		pair := [2]int{mv.x, mv.y}
		if used[pair] {
			return fmt.Errorf("move %d repeats overtake (%d,%d)", idx+1, mv.x, mv.y)
		}
		used[pair] = true

		px := pos[mv.x]
		py := pos[mv.y]
		if px != py+1 {
			return fmt.Errorf("move %d: car %d is not directly behind car %d", idx+1, mv.x, mv.y)
		}
		current[py], current[px] = current[px], current[py]
		pos[mv.x], pos[mv.y] = pos[mv.y], pos[mv.x]
	}
	for i := 0; i < n; i++ {
		if current[i] != tc.per[i] {
			return fmt.Errorf("final order mismatch at position %d: got %d expected %d", i+1, current[i], tc.per[i])
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sum := 0

	add := func(per []int) {
		if sum+len(per) > totalLimit {
			return
		}
		tests = append(tests, testCase{n: len(per), per: per})
		sum += len(per)
	}

	add([]int{2, 3, 1})
	add([]int{1})
	add([]int{1, 2})
	add([]int{3, 2, 1})

	for attempts := 0; attempts < 1000 && sum < totalLimit; attempts++ {
		n := rng.Intn(maxN) + 1
		if sum+n > totalLimit {
			n = totalLimit - sum
		}
		if n <= 0 {
			break
		}
		add(randPerm(n, rng))
	}

	if len(tests) == 0 {
		add([]int{1})
	}
	return tests
}

func buildInput(tc testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", tc.n)
	for i, v := range tc.per {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return b.String()
}

func randPerm(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
