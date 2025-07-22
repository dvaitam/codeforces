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

const mod = 1000000007

func solve(input string) string {
	rdr := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(rdr, &n); err != nil {
		return ""
	}
	a := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	B := 0
	for (1 << B) <= maxA {
		B++
	}
	size := 1 << B
	f := make([]int, size)
	for _, v := range a {
		f[v]++
	}
	for i := 0; i < B; i++ {
		bit := 1 << i
		for mask := 0; mask < size; mask++ {
			if mask&bit == 0 {
				f[mask] += f[mask|bit]
			}
		}
	}
	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	dp := make([]int, size)
	for mask := 0; mask < size; mask++ {
		cnt := f[mask]
		if cnt > 0 {
			dp[mask] = pow2[cnt] - 1
			if dp[mask] < 0 {
				dp[mask] += mod
			}
		}
	}
	for i := 0; i < B; i++ {
		bit := 1 << i
		for mask := 0; mask < size; mask++ {
			if mask&bit == 0 {
				dp[mask] -= dp[mask|bit]
				if dp[mask] < 0 {
					dp[mask] += mod
				}
			}
		}
	}
	return fmt.Sprintf("%d", dp[0])
}

func generateTests() []test {
	rand.Seed(452)
	var tests []test
	fixed := []string{
		"1\n0\n",
		"2\n1 1\n",
		"3\n1 2 3\n",
		"4\n0 0 0 0\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(10) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			val := rand.Intn(64)
			fmt.Fprintf(&sb, "%d", val)
			if i+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

type test struct {
	input, expected string
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
