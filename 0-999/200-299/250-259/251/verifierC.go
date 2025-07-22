package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func computeC(a, b int64, k int) int64 {
	diff := a - b
	if diff <= 0 {
		return 0
	}
	L := 1
	for x := 2; x <= k; x++ {
		L = lcm(L, x)
	}
	bmod := make([]int64, k+1)
	for x := 2; x <= k; x++ {
		bmod[x] = b % int64(x)
	}
	maxI := L
	if diff < int64(L) {
		maxI = int(diff)
	}
	dp := make([]int64, maxI+1)
	for i := 1; i <= maxI; i++ {
		dp[i] = dp[i-1] + 1
		for x := 2; x <= k; x++ {
			r := (bmod[x] + int64(i%x)) % int64(x)
			if r > 0 && int(r) <= i {
				cand := dp[i-int(r)] + 1
				if cand < dp[i] {
					dp[i] = cand
				}
			}
		}
	}
	if diff <= int64(maxI) {
		return dp[diff]
	}
	chunks := diff / int64(L)
	rem := diff % int64(L)
	return chunks*dp[L] + dp[rem]
}

type testCaseC struct {
	a, b int64
	k    int
	ans  int64
}

func genCaseC() testCaseC {
	rand.Seed(time.Now().UnixNano())
	b := rand.Int63n(50) + 1
	diff := rand.Int63n(60)
	a := b + diff
	k := rand.Intn(14) + 2
	return testCaseC{a, b, k, computeC(a, b, k)}
}

func buildInputC(cs []testCaseC) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintf(&sb, "%d %d %d\n", c.a, c.b, c.k)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCaseC, 100)
	for i := range cases {
		cases[i] = genCaseC()
	}
	input := buildInputC(cases)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cases) {
		fmt.Printf("expected %d lines got %d\n", len(cases), len(outputs))
		os.Exit(1)
	}
	for i, s := range outputs {
		var val int64
		fmt.Sscan(s, &val)
		if val != cases[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cases[i].ans, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
