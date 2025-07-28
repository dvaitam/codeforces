package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type interval struct {
	high int64
	low  int64
}

func find(parent map[int64]int64, x int64) int64 {
	if x <= 0 {
		return 0
	}
	root := x
	for {
		p, ok := parent[root]
		if !ok {
			parent[root] = root
			break
		}
		if p == root {
			break
		}
		root = p
	}
	for x != root {
		p := parent[x]
		parent[x] = root
		x = p
	}
	return root
}

func solveC(n int, a []int64) []int64 {
	intervals := make([]interval, n)
	for i := 0; i < n; i++ {
		intervals[i] = interval{high: a[i] + int64(i+1), low: a[i] + 1}
	}
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i].low == intervals[j].low {
			return intervals[i].high > intervals[j].high
		}
		return intervals[i].low > intervals[j].low
	})
	parent := make(map[int64]int64, n*2)
	res := make([]int64, 0, n)
	for _, it := range intervals {
		x := find(parent, it.high)
		if x >= it.low {
			parent[x] = find(parent, x-1)
			res = append(res, x)
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i] > res[j] })
	return res
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(20))
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	input := sb.String()
	ans := solveC(n, a)
	var exp strings.Builder
	for i, v := range ans {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(strconv.FormatInt(v, 10))
	}
	exp.WriteByte('\n')
	return input, exp.String()
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
