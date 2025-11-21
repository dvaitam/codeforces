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

const refSource = "0-999/300-399/330-339/335/335C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns := strings.TrimSpace(refOut)

		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		ans := strings.TrimSpace(candOut)
		if ans != "WIN" && ans != "LOSE" {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid output %q\n", idx+1, tc.name, ans)
			os.Exit(1)
		}
		if ans != refAns {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, refAns, ans)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier335C-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "empty", input: "3 0\n"},
		{name: "sample1", input: "3 1\n1 1\n"},
		{name: "sample2", input: "12 2\n4 1\n8 1\n"},
		{name: "sample3", input: "1 1\n1 2\n"},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, fmt.Sprintf("random_%d", i+1)))
	}
	return tests
}

func randomTest(rng *rand.Rand, name string) testCase {
	r := rng.Intn(100) + 1
	states := make([]int, r)
	for i := range states {
		states[i] = 3
	}
	type cell struct{ row, col int }
	var chosen []cell
	available := collectMoves(states)
	maxMoves := rng.Intn(r + 1)
	for move := 0; move < maxMoves && len(available) > 0; move++ {
		c := available[rng.Intn(len(available))]
		applyMove(states, c.row, c.col)
		chosen = append(chosen, c)
		available = collectMoves(states)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", r, len(chosen))
	for _, c := range chosen {
		fmt.Fprintf(&sb, "%d %d\n", c.row+1, c.col+1)
	}
	return testCase{name: name, input: sb.String()}
}

func collectMoves(states []int) []struct{ row, col int } {
	var moves []struct{ row, col int }
	for i, mask := range states {
		if mask&1 != 0 {
			moves = append(moves, struct{ row, col int }{i, 0})
		}
		if mask&2 != 0 {
			moves = append(moves, struct{ row, col int }{i, 1})
		}
	}
	return moves
}

func applyMove(states []int, row int, col int) {
	if row < 0 || row >= len(states) {
		return
	}
	states[row] &^= 1 << col
	blockBit := 1 << (1 - col)
	for delta := -1; delta <= 1; delta++ {
		if delta == 0 {
			continue
		}
		nr := row + delta
		if nr >= 0 && nr < len(states) {
			states[nr] &^= blockBit
		}
	}
}
