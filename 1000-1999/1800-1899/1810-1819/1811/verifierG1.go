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

func solve(tc testG1) int {
	n := tc.n
	k := tc.k
	// 1-based indexing for logic matching
	a := make([]int, n+1)
	for i := 0; i < n; i++ {
		a[i+1] = tc.colors[i]
	}

	// Combinations table
	C := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		C[i] = make([]int, n+1)
		C[i][0] = 1
		if i > 0 {
			for j := 1; j < i; j++ {
				C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
			}
			C[i][i] = 1
		}
	}

	f := make([]int, n+1)
	g := make([]int, n+1)
	g[0] = 1

	for i := 1; i <= n; i++ {
		s := 0
		for j := i - 1; j >= 0; j-- {
			if s >= k-1 {
				ways := C[s][k-1]
				if f[j]+1 > f[i] {
					f[i] = f[j] + 1
					g[i] = (g[j] * ways) % MOD
				} else if f[j]+1 == f[i] {
					g[i] = (g[i] + g[j]*ways) % MOD
				}
			}
			if j > 0 && a[j] == a[i] {
				s++
			}
		}
	}

	mx := 0
	ans := 1
	for i := 1; i <= n; i++ {
		if f[i] > mx {
			mx = f[i]
			ans = g[i]
		} else if f[i] == mx {
			ans = (ans + g[i]) % MOD
		}
	}
	return ans
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
	fmt.Println("All tests passed")
}