package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input  string
	output string
}

func simulate(r, b int, start byte) int {
	rRem, bRem := r, b
	last := start
	if start == 'R' {
		rRem--
	} else {
		bRem--
	}
	blocks := 1
	total := r + b
	for move := 2; move <= total; move++ {
		if move%2 == 1 {
			if last == 'R' {
				if rRem > 0 {
					rRem--
				} else {
					bRem--
					last = 'B'
					blocks++
				}
			} else {
				if bRem > 0 {
					bRem--
				} else {
					rRem--
					last = 'R'
					blocks++
				}
			}
		} else {
			if last == 'R' {
				if bRem > 0 {
					bRem--
					last = 'B'
					blocks++
				} else {
					rRem--
				}
			} else {
				if rRem > 0 {
					rRem--
					last = 'R'
					blocks++
				} else {
					bRem--
				}
			}
		}
	}
	return blocks
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(in, &n, &m)
	blocksR := simulate(n, m, 'R')
	blocksB := simulate(n, m, 'B')
	blocks := blocksR
	if blocksB < blocksR {
		blocks = blocksB
	}
	total := n + m
	petya := total - blocks
	vasya := blocks - 1
	return fmt.Sprintf("%d %d\n", petya, vasya)
}

func generateTests() []testCase {
	rand.Seed(43)
	var tests []testCase
	fixed := []string{
		"1 1\n",
		"2 1\n",
		"1 2\n",
		"3 3\n",
	}
	for _, f := range fixed {
		tests = append(tests, testCase{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		inp := fmt.Sprintf("%d %d\n", n, m)
		tests = append(tests, testCase{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %sGot: %s\n", i+1, t.input, strings.TrimSpace(t.output), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
