package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

type testCaseE struct {
	n int
	k int
}

func genTestsE() []testCaseE {
	rand.Seed(5)
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rand.Intn(8) + 1 //1..8
		k := rand.Intn(n + 1)
		tests[i] = testCaseE{n, k}
	}
	return tests
}

func solveE(tc testCaseE) int {
	n, k := tc.n, tc.k
	if k > n {
		return 0
	}
	maxExp := n * n
	pow2 := make([]int, maxExp+1)
	pow2[0] = 1
	for i := 1; i <= maxExp; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}
	A := make([][]int, n+1)
	for prev := 0; prev <= n; prev++ {
		A[prev] = make([]int, n+1)
		for i := prev + 1; i <= n; i++ {
			d := i - prev
			v := pow2[d] - 1
			if v < 0 {
				v += mod
			}
			e := d * prev
			v = int(int64(v) * int64(pow2[e]) % mod)
			A[prev][i] = v
		}
	}
	dp := make([][]int, k+1)
	for t := 0; t <= k; t++ {
		dp[t] = make([]int, n+1)
	}
	dp[0][0] = 1
	for t := 1; t <= k; t++ {
		for i := 1; i <= n; i++ {
			var sum int64
			for prev := 0; prev < i; prev++ {
				if dp[t-1][prev] != 0 {
					sum += int64(dp[t-1][prev]) * int64(A[prev][i])
					sum %= mod
				}
			}
			dp[t][i] = int(sum % mod)
		}
	}
	var ans int64
	for i := 0; i <= n; i++ {
		if dp[k][i] == 0 {
			continue
		}
		e := (n - i) * i
		ans = (ans + int64(dp[k][i])*int64(pow2[e])) % mod
	}
	return int(ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		expect := solveE(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
