package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

const MOD int = 1_000_000_007
const NEG_INF int = -1_000_000_000

type testG1 struct {
	n      int
	k      int
	colors []int
}

func genTests() []testG1 {
	rand.Seed(181107)
	tests := make([]testG1, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		colors := make([]int, n)
		for j := range colors {
			colors[j] = rand.Intn(n) + 1
		}
		tests[i] = testG1{n: n, k: k, colors: colors}
	}
	return tests
}

type state struct {
	len int
	cnt int
}

func solve(tc testG1) int {
	n := tc.n
	k := tc.k
	colors := tc.colors
	dpPrev := make([][]state, k)
	dpCur := make([][]state, k)
	for r := 0; r < k; r++ {
		dpPrev[r] = make([]state, n+1)
		dpCur[r] = make([]state, n+1)
		for c := 0; c <= n; c++ {
			dpPrev[r][c].len = NEG_INF
			dpCur[r][c].len = NEG_INF
		}
	}
	dpPrev[0][0] = state{0, 1}

	for idx := n - 1; idx >= 0; idx-- {
		for r := 0; r < k; r++ {
			for c := 0; c <= n; c++ {
				dpCur[r][c].len = NEG_INF
				dpCur[r][c].cnt = 0
			}
		}
		color := colors[idx]
		for r := 0; r < k; r++ {
			for c := 0; c <= n; c++ {
				best := dpPrev[r][c]
				if best.len == NEG_INF {
					dpCur[r][c] = best
					continue
				}
				curLen := best.len
				curCnt := best.cnt
				if r == 0 {
					nr := 1
					nc := color
					if nr == k {
						nr = 0
						nc = 0
					}
					cand := dpPrev[nr][nc]
					if cand.len != NEG_INF {
						candLen := cand.len + 1
						candCnt := cand.cnt
						if candLen > curLen {
							curLen = candLen
							curCnt = candCnt
						} else if candLen == curLen {
							curCnt = (curCnt + candCnt) % MOD
						}
					}
				} else if c == color {
					nr := r + 1
					nc := c
					if nr == k {
						nr = 0
						nc = 0
					}
					cand := dpPrev[nr][nc]
					if cand.len != NEG_INF {
						candLen := cand.len + 1
						candCnt := cand.cnt
						if candLen > curLen {
							curLen = candLen
							curCnt = candCnt
						} else if candLen == curLen {
							curCnt = (curCnt + candCnt) % MOD
						}
					}
				}
				dpCur[r][c].len = curLen
				dpCur[r][c].cnt = curCnt % MOD
			}
		}
		dpPrev, dpCur = dpCur, dpPrev
	}
	return dpPrev[0][0].cnt % MOD
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.colors {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}

	expected := make([]int, len(tests))
	for i, tc := range tests {
		expected[i] = solve(tc)
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
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
