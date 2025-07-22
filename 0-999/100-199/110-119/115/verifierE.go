package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type race struct {
	l, r int
	p    int64
}

func solveE(n int, costs []int64, races []race) int64 {
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + costs[i-1]
	}
	byL := make([][]race, n+2)
	for _, rc := range races {
		if rc.l >= 1 && rc.l <= n {
			byL[rc.l] = append(byL[rc.l], rc)
		}
	}
	dp := make([]int64, n+3)
	for i := n; i >= 1; i-- {
		best := dp[i+1]
		for _, rc := range byL[i] {
			costSeg := prefix[rc.r] - prefix[i-1]
			temp := dp[rc.r+1] + rc.p - costSeg
			if temp > best {
				best = temp
			}
		}
		dp[i] = best
	}
	if dp[1] < 0 {
		return 0
	}
	return dp[1]
}

func genCase() (string, string) {
	n := rand.Intn(10) + 1
	m := rand.Intn(10)
	costs := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		costs[i] = int64(rand.Intn(10))
		fmt.Fprintf(&sb, "%d\n", costs[i])
	}
	races := make([]race, m)
	for i := 0; i < m; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n-l+1) + l
		p := int64(rand.Intn(20))
		races[i] = race{l, l, 0}
		races[i].r = r
		races[i].p = p
		fmt.Fprintf(&sb, "%d %d %d\n", l, r, p)
	}
	input := sb.String()
	expect := solveE(n, costs, races)
	return input, fmt.Sprint(expect)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		in, expect := genCase()
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			fmt.Println(in)
			return
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %s\nGot: %s\n", t, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
