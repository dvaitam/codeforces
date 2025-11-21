package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()

	for idx, s := range tests {
		expected := countPreimages(s)
		candOut, err := runProgram(candidate, s+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			os.Exit(1)
		}
		got, err := parseBigInt(strings.TrimSpace(candOut))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			os.Exit(1)
		}
		if expected.Cmp(got) != 0 {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %s, got %s\nInput:\n%s\n", idx+1, expected.String(), got.String(), s)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []string {
	tests := []string{
		"BABBBABBA",
		"ABABB",
		"ABABAB",
		"AAA",
		"BBB",
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		n := rng.Intn(98) + 3
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb[i] = 'A'
			} else {
				sb[i] = 'B'
			}
		}
		tests = append(tests, string(sb))
	}
	return tests
}

func countPreimages(s string) *big.Int {
	n := len(s)
	target := make([]int, n)
	for i := 0; i < n; i++ {
		if s[i] == 'B' {
			target[i] = 1
		} else {
			target[i] = 0
		}
	}
	total := new(big.Int)
	for t0 := 0; t0 < 2; t0++ {
		for t1 := 0; t1 < 2; t1++ {
			states := make([]*big.Int, 4)
			states[encodeState(t0, t1)] = big.NewInt(1)
			for i := 1; i <= n-2; i++ {
				nextStates := make([]*big.Int, 4)
				for st, val := range states {
					if val == nil {
						continue
					}
					prev := st & 1
					cur := (st >> 1) & 1
					for next := 0; next < 2; next++ {
						if produce(prev, cur, next) == target[i] {
							idx := encodeState(cur, next)
							if nextStates[idx] == nil {
								nextStates[idx] = new(big.Int)
							}
							nextStates[idx].Add(nextStates[idx], val)
						}
					}
				}
				states = nextStates
			}
			for st, val := range states {
				if val == nil {
					continue
				}
				prev := st & 1
				cur := (st >> 1) & 1
				if produce(prev, cur, t0) != target[n-1] {
					continue
				}
				if produce(cur, t0, t1) != target[0] {
					continue
				}
				total.Add(total, val)
			}
		}
	}
	return total
}

func produce(left, cur, right int) int {
	if cur == 0 && right == 1 {
		return 1 // AB -> BA, so current becomes B
	}
	if left == 0 && cur == 1 {
		return 0 // previous pair AB affected current B -> A
	}
	return cur
}

func encodeState(prev, cur int) int {
	return prev | (cur << 1)
}

func parseBigInt(s string) (*big.Int, error) {
	if s == "" {
		return nil, fmt.Errorf("empty output")
	}
	val, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer %q", s)
	}
	return val, nil
}
