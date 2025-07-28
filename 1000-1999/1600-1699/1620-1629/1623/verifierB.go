package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type interval struct {
	l int
	r int
}

type intervalAns struct {
	l int
	r int
	d int
}

func expected(n int, intervals []interval) string {
	ivs := make([]intervalAns, n)
	for i, v := range intervals {
		ivs[i] = intervalAns{l: v.l, r: v.r}
	}
	sort.Slice(ivs, func(i, j int) bool {
		return (ivs[i].r - ivs[i].l) < (ivs[j].r - ivs[j].l)
	})
	flag := make([]bool, n+2)
	for i := 0; i < n; i++ {
		for k := ivs[i].l; k <= ivs[i].r; k++ {
			if k >= 0 && k < len(flag) && !flag[k] {
				ivs[i].d = k
				flag[k] = true
				break
			}
		}
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", ivs[i].l, ivs[i].r, ivs[i].d))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateIntervals(n int, rng *rand.Rand) []interval {
	intervals := make([]interval, 0, n)
	var dfs func(l, r int)
	dfs = func(l, r int) {
		if l > r {
			return
		}
		d := rng.Intn(r-l+1) + l
		intervals = append(intervals, interval{l: l, r: r})
		if l <= d-1 {
			dfs(l, d-1)
		}
		if d+1 <= r {
			dfs(d+1, r)
		}
	}
	dfs(1, n)
	rng.Shuffle(len(intervals), func(i, j int) { intervals[i], intervals[j] = intervals[j], intervals[i] })
	return intervals
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(7) + 1
		intervals := generateIntervals(n, rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, iv := range intervals {
			sb.WriteString(fmt.Sprintf("%d %d\n", iv.l, iv.r))
		}
		input := sb.String()
		exp := strings.TrimSpace(expected(n, intervals))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\nGot:\n%s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
