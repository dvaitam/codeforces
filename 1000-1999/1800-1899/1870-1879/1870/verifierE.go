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

const (
	MaxXor = 1 << 13
	Words  = MaxXor / 64
	Limit  = 130
)

type Bitset [Words]uint64

type testCaseE struct {
	n int
	a []int
}

func applyXor(dst *Bitset, src *Bitset, val int) {
	for i := 0; i < Words; i++ {
		w := src[i]
		for w != 0 {
			b := bits.TrailingZeros64(w)
			idx := i*64 + b
			nidx := idx ^ val
			dst[nidx>>6] |= 1 << (uint(nidx) & 63)
			w &= w - 1
		}
	}
}

func genTestsE() []testCaseE {
	rng := rand.New(rand.NewSource(46))
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rng.Intn(6) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rng.Intn(10)
		}
		tests[i] = testCaseE{n, a}
	}
	return tests
}

func solveE(tc testCaseE) int {
	n := tc.n
	a := tc.a
	dp := make([]Bitset, n+1)
	dp[0][0] = 1
	freq := make([]int, Limit+2)
	for r := 1; r <= n; r++ {
		for i := 0; i < Words; i++ {
			dp[r][i] = dp[r-1][i]
		}
		for i := range freq {
			freq[i] = 0
		}
		mex := 0
		for l := r; l >= 1 && r-l+1 <= Limit; l-- {
			v := a[l-1]
			if v <= Limit {
				freq[v]++
			}
			for mex <= Limit && freq[mex] > 0 {
				mex++
			}
			if mex > Limit {
				break
			}
			applyXor(&dp[r], &dp[l-1], mex)
		}
	}
	ans := 0
	for x := MaxXor - 1; x >= 0; x-- {
		if (dp[n][x>>6]>>(uint(x)&63))&1 != 0 {
			ans = x
			break
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
	}
	expected := make([]int, len(tests))
	for i, tc := range tests {
		expected[i] = solveE(tc)
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
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
