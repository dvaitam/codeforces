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

func bruteProfit(c, s, d, l []int) int {
	n := len(c)
	m := len(d)
	type state struct{ i, mask int }
	memo := make(map[state]int)
	var dfs func(i, mask int) int
	dfs = func(i, mask int) int {
		if i == n {
			return 0
		}
		st := state{i, mask}
		if val, ok := memo[st]; ok {
			return val
		}
		res := dfs(i+1, mask)
		for j := 0; j < m; j++ {
			if mask&(1<<j) == 0 && c[i] <= d[j] && (l[j] == s[i] || l[j] == s[i]-1) {
				val := c[i] + dfs(i+1, mask|(1<<j))
				if val > res {
					res = val
				}
			}
		}
		memo[st] = res
		return res
	}
	return dfs(0, 0)
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 2
	c := make([]int, n)
	s := make([]int, n)
	usedSizes := rand.Perm(50)
	var sb strings.Builder
	fmt.Fprintln(&sb, n)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(20) + 1
		s[i] = usedSizes[i] + 1
		fmt.Fprintf(&sb, "%d %d\n", c[i], s[i])
	}
	m := rng.Intn(4) + 2
	dVal := make([]int, m)
	lVal := make([]int, m)
	fmt.Fprintln(&sb, m)
	for i := 0; i < m; i++ {
		dVal[i] = rng.Intn(20) + 1
		lVal[i] = rng.Intn(50) + 1
		fmt.Fprintf(&sb, "%d %d\n", dVal[i], lVal[i])
	}
	exp := bruteProfit(c, s, dVal, lVal)
	return sb.String(), exp
}

func runCase(exe, input string, expectedAns int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expectedAns {
		return fmt.Errorf("expected %d got %d", expectedAns, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
