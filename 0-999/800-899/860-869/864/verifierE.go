package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct {
	profit int64
	used   []bool
}

func solveE(tt []int, dd []int, pp []int64) (int64, []int) {
	n := len(tt)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return dd[idx[i]] < dd[idx[j]] })
	maxD := 0
	for _, d := range dd {
		if d > maxD {
			maxD = d
		}
	}
	dp := make([]pair, maxD+1)
	for t := range dp {
		dp[t].used = make([]bool, n)
	}
	best := pair{used: make([]bool, n)}
	for _, x := range idx {
		t := tt[x]
		d := dd[x]
		p := pp[x]
		for tim := d - t - 1; tim >= 0; tim-- {
			cur := dp[tim]
			newProfit := cur.profit + p
			newT := tim + t
			if newProfit > dp[newT].profit {
				newUsed := make([]bool, n)
				copy(newUsed, cur.used)
				newUsed[x] = true
				dp[newT].profit = newProfit
				dp[newT].used = newUsed
				if newProfit > best.profit {
					best.profit = newProfit
					best.used = newUsed
				}
			}
		}
	}
	list := []int{}
	for i := 0; i < n; i++ {
		if best.used[i] {
			list = append(list, i)
		}
	}
	sort.Slice(list, func(i, j int) bool { return dd[list[i]] < dd[list[j]] })
	return best.profit, list
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	tt := make([]int, n)
	dd := make([]int, n)
	pp := make([]int64, n)
	for i := 0; i < n; i++ {
		tt[i] = rng.Intn(20) + 1
		dd[i] = rng.Intn(200) + tt[i] + 1
		pp[i] = int64(rng.Intn(20) + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tt[i], dd[i], pp[i]))
	}
	profit, list := solveE(tt, dd, pp)
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%d\n", profit))
	out.WriteString(fmt.Sprintf("%d\n", len(list)))
	for _, v := range list {
		out.WriteString(fmt.Sprintf("%d ", v+1))
	}
	return sb.String(), strings.TrimSpace(out.String())
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
