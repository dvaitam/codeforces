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

func calcWays(k int, a, b int) int {
	const mod = 1000000007
	dp := make([]int, b+1)
	dp[0] = 1
	for i := 1; i <= b; i++ {
		dp[i] = dp[i-1]
		if i >= k {
			dp[i] = (dp[i] + dp[i-k]) % mod
		}
	}
	pref := make([]int, b+1)
	for i := 1; i <= b; i++ {
		pref[i] = pref[i-1] + dp[i]
		if pref[i] >= mod {
			pref[i] -= mod
		}
	}
	res := pref[b]
	if a > 1 {
		res -= pref[a-1]
	}
	if res < 0 {
		res += mod
	}
	return res
}

func genCase(rng *rand.Rand) (int, int, [][2]int) {
	t := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	pairs := make([][2]int, t)
	for i := range pairs {
		a := rng.Intn(20) + 1
		b := a + rng.Intn(20)
		pairs[i] = [2]int{a, b}
	}
	return k, t, pairs
}

func runCase(bin string, k, t int, pairs [][2]int, exp []int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", t, k)
	for i := range pairs {
		fmt.Fprintf(&sb, "%d %d\n", pairs[i][0], pairs[i][1])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Fields(strings.TrimSpace(out.String()))
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, l := range lines {
		var val int
		fmt.Sscan(l, &val)
		if val != exp[i] {
			return fmt.Errorf("test %d expected %d got %d", i+1, exp[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		k, t, pairs := genCase(rng)
		exp := make([]int, t)
		for j := 0; j < t; j++ {
			exp[j] = calcWays(k, pairs[j][0], pairs[j][1])
		}
		if err := runCase(bin, k, t, pairs, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
