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

const MOD int = 1000000007

type testG2 struct {
	n      int
	k      int
	colors []int
}

func genTests() []testG2 {
	rand.Seed(181108)
	tests := make([]testG2, 100)
	for i := range tests {
		n := rand.Intn(20) + 1
		k := rand.Intn(n) + 1
		colors := make([]int, n)
		for j := range colors {
			colors[j] = rand.Intn(n) + 1
		}
		tests[i] = testG2{n: n, k: k, colors: colors}
	}
	return tests
}

func solveCase(n, k int, colors []int) int {
	freq := make([]int, n+1)
	for _, c := range colors {
		freq[c]++
	}
	lens := make([][]int, n+1)
	cnts := make([][]int, n+1)
	for c := 1; c <= n; c++ {
		if freq[c] > 0 {
			size := k
			if freq[c] < size {
				size = freq[c]
			}
			lens[c] = make([]int, size+1)
			cnts[c] = make([]int, size+1)
			for i := 1; i <= size; i++ {
				lens[c][i] = -1
			}
		}
	}
	bestLen := 0
	bestCnt := 1
	for _, col := range colors {
		arrLen := lens[col]
		arrCnt := cnts[col]
		size := len(arrLen) - 1
		oldKLen, oldKCnt := 0, 0
		if k <= size {
			oldKLen = arrLen[k]
			oldKCnt = arrCnt[k]
		}
		for r := size; r >= 1; r-- {
			baseLen := -1
			baseCnt := 0
			if r == 1 {
				baseLen = bestLen
				baseCnt = bestCnt
			} else if arrLen[r-1] != -1 {
				baseLen = arrLen[r-1]
				baseCnt = arrCnt[r-1]
			}
			if baseLen != -1 {
				newLen := baseLen + 1
				if newLen > arrLen[r] {
					arrLen[r] = newLen
					arrCnt[r] = baseCnt
				} else if newLen == arrLen[r] {
					arrCnt[r] += baseCnt
					if arrCnt[r] >= MOD {
						arrCnt[r] -= MOD
					}
				}
			}
		}
		if k <= size {
			if arrLen[k] > bestLen {
				bestLen = arrLen[k]
				bestCnt = arrCnt[k] % MOD
			} else if arrLen[k] == bestLen {
				added := arrCnt[k]
				if oldKLen == bestLen {
					added -= oldKCnt
					added %= MOD
					if added < 0 {
						added += MOD
					}
				}
				bestCnt += added
				bestCnt %= MOD
			}
		}
	}
	return bestCnt % MOD
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG2.go /path/to/binary")
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
		expected[i] = solveCase(tc.n, tc.k, tc.colors)
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
