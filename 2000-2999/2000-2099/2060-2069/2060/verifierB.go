package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceB = "2000-2999/2000-2099/2060-2069/2060/2060B.go"

type testCase struct {
	n, m      int
	cards     [][]int
	possible  bool
	residueOf []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	prepareTests(tests)

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	if _, err := runCommand(exec.Command(refBin), input); err != nil {
		// Reference is trusted; still show its stderr for debugging.
		fail("reference execution failed: %v", err)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	userTokens := strings.Fields(userOut)

	idx := 0
	for ti, tc := range tests {
		if !tc.possible {
			if idx >= len(userTokens) {
				fail("test %d: missing output", ti+1)
			}
			if userTokens[idx] != "-1" {
				fail("test %d: expected -1, got %q", ti+1, userTokens[idx])
			}
			idx++
			continue
		}
		if idx+tc.n > len(userTokens) {
			fail("test %d: expected %d numbers, got %d", ti+1, tc.n, len(userTokens)-idx)
		}
		perm := make([]int, tc.n)
		seen := make([]bool, tc.n+1)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(userTokens[idx+i])
			if err != nil {
				fail("test %d: invalid integer %q", ti+1, userTokens[idx+i])
			}
			if val < 1 || val > tc.n || seen[val] {
				fail("test %d: output is not a permutation", ti+1)
			}
			seen[val] = true
			perm[i] = val - 1 // zero-based cow index
		}
		idx += tc.n

		for pos, cow := range perm {
			if tc.residueOf[cow] != pos {
				fail("test %d: cow at position %d has residue %d, expected %d", ti+1, pos+1, tc.residueOf[cow], pos)
			}
		}
	}
	if idx != len(userTokens) {
		fail("extra tokens in output")
	}

	fmt.Println("OK")
}

func parseInput(data []byte) ([]testCase, error) {
	in := bufio.NewReader(strings.NewReader(string(data)))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return nil, err
		}
		cards := make([][]int, n)
		for j := 0; j < n; j++ {
			cards[j] = make([]int, m)
			for k := 0; k < m; k++ {
				if _, err := fmt.Fscan(in, &cards[j][k]); err != nil {
					return nil, err
				}
			}
		}
		tests[i] = testCase{n: n, m: m, cards: cards}
	}
	return tests, nil
}

func prepareTests(tests []testCase) {
	for idx := range tests {
		tc := &tests[idx]
		resPos := make([]int, tc.n)
		for i := range resPos {
			resPos[i] = -1
		}
		tc.residueOf = make([]int, tc.n)
		valid := true
		for cow := 0; cow < tc.n && valid; cow++ {
			res := tc.cards[cow][0] % tc.n
			seen := make([]bool, tc.m)
			for _, v := range tc.cards[cow] {
				if v%tc.n != res {
					valid = false
					break
				}
				q := v / tc.n
				if q < 0 || q >= tc.m || seen[q] {
					valid = false
					break
				}
				seen[q] = true
			}
			if !valid {
				break
			}
			for _, s := range seen {
				if !s {
					valid = false
					break
				}
			}
			if !valid {
				break
			}
			if resPos[res] != -1 {
				valid = false
				break
			}
			resPos[res] = cow
			tc.residueOf[cow] = res
		}
		for _, pos := range resPos {
			if pos == -1 {
				valid = false
				break
			}
		}
		tc.possible = valid
	}
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2060B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceB))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
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

func runCommand(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
