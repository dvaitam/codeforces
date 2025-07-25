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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Fenwick struct {
	n int
	t []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, t: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
	for i <= f.n {
		f.t[i] += v
		i += i & -i
	}
}

func (f *Fenwick) Sum(i int) int {
	s := 0
	for i > 0 {
		s += f.t[i]
		i -= i & -i
	}
	return s
}

func expected(n int, R int64, ducks [][2]int64) string {
	type interval struct{ l, r int64 }
	intervals := make([]interval, 0, n)
	for i := 0; i < n; i++ {
		hi, ti := ducks[i][0], ducks[i][1]
		if ti < 0 {
			continue
		}
		l := hi
		if l < 0 {
			l = 0
		}
		intervals = append(intervals, interval{l, ti})
	}
	if len(intervals) == 0 {
		return "0"
	}
	sort.Slice(intervals, func(i, j int) bool { return intervals[i].r < intervals[j].r })
	m := len(intervals)
	rArr := make([]int64, m+1)
	llist := make([]struct {
		l   int64
		idx int
	}, m)
	for i := 1; i <= m; i++ {
		rArr[i] = intervals[i-1].r
		llist[i-1] = struct {
			l   int64
			idx int
		}{intervals[i-1].l, i}
	}
	sort.Slice(llist, func(i, j int) bool { return llist[i].l < llist[j].l })
	fw := NewFenwick(m)
	count := make([]int, m+1)
	ptr := 0
	for i := 1; i <= m; i++ {
		ri := rArr[i]
		for ptr < m && llist[ptr].l <= ri {
			fw.Add(llist[ptr].idx, 1)
			ptr++
		}
		count[i] = fw.Sum(m) - fw.Sum(i-1)
	}
	p := make([]int, m+1)
	for i := 1; i <= m; i++ {
		val := rArr[i] - R
		k := sort.Search(m, func(j int) bool { return rArr[j+1] > val })
		p[i] = k
	}
	dp := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		without := dp[i-1]
		with := dp[p[i]] + int64(count[i])
		if without > with {
			dp[i] = without
		} else {
			dp[i] = with
		}
	}
	return fmt.Sprintf("%d", dp[m])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		R := int64(rng.Intn(5) + 1)
		ducks := make([][2]int64, n)
		for j := 0; j < n; j++ {
			h := int64(rng.Intn(11) - 5)
			t := h + int64(rng.Intn(6)+1)
			ducks[j] = [2]int64{h, t}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, R))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", ducks[j][0], ducks[j][1]))
		}
		input := sb.String()
		exp := expected(n, R, ducks)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
