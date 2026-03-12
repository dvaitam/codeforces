package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod int64 = 998244353

func solveCase(n int64, k int) int64 {
	if k == 0 {
		return 0
	}
	f := make([][]int64, n+2)
	s := make([][]int64, n+2)
	for i := range f {
		f[i] = make([]int64, k+1)
		s[i] = make([]int64, k+1)
	}

	var ans int64
	for i := int64(1); i <= n; i++ {
		f[i][1] = 1
		s[i][1] = i % mod
		if i > 1 {
			for j := 2; j <= k; j++ {
				term1 := (s[i-1][j-1] - f[i-2][j-1] + mod) % mod
				term2 := (s[i-1][j] - f[i-2][j] + mod) % mod
				if j == k {
					f[i][j] = (term1 + term2) % mod
				} else {
					f[i][j] = term1 % mod
				}
				s[i][j] = (s[i-1][j] + f[i][j]) % mod
			}
		}
		if i < n && k >= 2 {
			ans = (ans + f[i][k] + f[i][k-1]) % mod
		} else if i < n && k == 1 {
			ans = (ans + f[i][k]) % mod
		}
	}
	return ans % mod
}

type testCase struct {
	n int64
	k int
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.n, tc.k)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}

	rnd := rand.New(rand.NewSource(1989))

	var tests []testCase

	// Small deterministic cases
	for n := int64(2); n <= 10; n++ {
		maxK := n
		if maxK > 10 {
			maxK = 10
		}
		for k := 2; int64(k) <= maxK; k++ {
			tests = append(tests, testCase{n: n, k: k})
		}
	}

	// Random medium cases
	for i := 0; i < 40; i++ {
		n := int64(rnd.Intn(198)) + 3 // 3..200
		maxK := n
		if maxK > 10 {
			maxK = 10
		}
		k := rnd.Intn(int(maxK)-1) + 2 // 2..min(n,10)
		tests = append(tests, testCase{n: n, k: k})
	}

	// A few larger cases
	for i := 0; i < 10; i++ {
		n := int64(rnd.Intn(1000)) + 100
		maxK := n
		if maxK > 10 {
			maxK = 10
		}
		k := rnd.Intn(int(maxK)-1) + 2
		tests = append(tests, testCase{n: n, k: k})
	}

	for i, tc := range tests {
		input := buildInput(tc)
		want := strconv.FormatInt(solveCase(tc.n, tc.k), 10)
		got, err := runCandidate(os.Args[1], input)
		if err != nil {
			fmt.Printf("case %d (n=%d k=%d) runtime error: %v\n", i+1, tc.n, tc.k, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != want {
			fmt.Printf("case %d (n=%d k=%d) failed\nexpected: %s\ngot: %s\n", i+1, tc.n, tc.k, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
