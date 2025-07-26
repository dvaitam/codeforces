package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type itemB struct {
	cost int
	mask int
}

type testCaseB struct {
	items []itemB
}

func solveB(tc testCaseB) int {
	const INF = int(1e9)
	dp := [8]int{}
	for i := 1; i < 8; i++ {
		dp[i] = INF
	}
	for _, it := range tc.items {
		for mask := 7; mask >= 0; mask-- {
			if dp[mask] == INF {
				continue
			}
			nm := mask | it.mask
			val := dp[mask] + it.cost
			if val < dp[nm] {
				dp[nm] = val
			}
		}
		if it.cost < dp[it.mask] {
			dp[it.mask] = it.cost
		}
	}
	if dp[7] >= INF {
		return -1
	}
	return dp[7]
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseB {
	rng := rand.New(rand.NewSource(43))
	tests := make([]testCaseB, 100)
	for i := range tests {
		n := rng.Intn(6) + 1
		items := make([]itemB, n)
		for j := range items {
			cost := rng.Intn(50) + 1
			mask := 0
			if rng.Intn(2) == 1 {
				mask |= 1
			}
			if rng.Intn(2) == 1 {
				mask |= 2
			}
			if rng.Intn(2) == 1 {
				mask |= 4
			}
			if mask == 0 {
				mask = 1 << uint(rng.Intn(3))
			}
			items[j] = itemB{cost, mask}
		}
		tests[i] = testCaseB{items}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(tc.items))
		for _, it := range tc.items {
			vitamins := ""
			if it.mask&1 != 0 {
				vitamins += "A"
			}
			if it.mask&2 != 0 {
				vitamins += "B"
			}
			if it.mask&4 != 0 {
				vitamins += "C"
			}
			fmt.Fprintf(&sb, "%d %s\n", it.cost, vitamins)
		}
		exp := solveB(tc)
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		var got int
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Printf("test %d failed:\ninput:%sexpected %d got %s\n", i+1, sb.String(), exp, out)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
