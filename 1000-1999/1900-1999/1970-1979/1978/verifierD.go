package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type testCaseD struct {
	n int
	c int
	a []int
}

type subset struct {
	mask int
}

func genTestsD() []testCaseD {
	rand.Seed(45)
	tests := make([]testCaseD, 20)
	for i := range tests {
		n := rand.Intn(5) + 2 // 2..6
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(5)
		}
		tests[i] = testCaseD{n: n, c: rand.Intn(5), a: a}
	}
	return tests
}

func winner(n int, c int, a []int, mask int) int {
	undecided := c
	remaining := make([]int, n)
	minIdx := -1
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			undecided += a[i]
		} else {
			remaining[i] = a[i]
			if minIdx == -1 {
				minIdx = i
			}
		}
	}
	if minIdx != -1 {
		remaining[minIdx] += undecided
	}
	winner := -1
	best := -1
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			continue
		}
		v := remaining[i]
		if v > best || (v == best && i < winner) {
			best = v
			winner = i
		}
	}
	return winner
}

func solveD(tc testCaseD) []int {
	res := make([]int, tc.n)
	for target := 0; target < tc.n; target++ {
		best := tc.n + 1
		for mask := 0; mask < 1<<tc.n; mask++ {
			if mask&(1<<target) != 0 {
				continue
			}
			if winner(tc.n, tc.c, tc.a, mask) == target {
				removed := bits.OnesCount(uint(mask))
				if removed < best {
					best = removed
				}
			}
		}
		res[target] = best
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.c)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([][]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveD(tc)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, tc := range tests {
		for j := 0; j < tc.n; j++ {
			if !scanner.Scan() {
				fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
				os.Exit(1)
			}
			val, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
				os.Exit(1)
			}
			if val != expected[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
				os.Exit(1)
			}
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
