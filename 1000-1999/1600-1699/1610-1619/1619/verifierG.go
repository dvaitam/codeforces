package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type dsu struct {
	parent []int
	rank   []byte
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), rank: make([]byte, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.rank[a] < d.rank[b] {
		a, b = b, a
	}
	d.parent[b] = a
	if d.rank[a] == d.rank[b] {
		d.rank[a]++
	}
}

type mine struct {
	x, y int
	t    int
}

func solveCase(mines []mine, k int) int {
	n := len(mines)
	d := newDSU(n)
	type pair struct{ val, idx int }
	mpX := make(map[int][]pair)
	for i, m := range mines {
		mpX[m.x] = append(mpX[m.x], pair{m.y, i})
	}
	for _, arr := range mpX {
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		for j := 1; j < len(arr); j++ {
			if arr[j].val-arr[j-1].val <= k {
				d.union(arr[j].idx, arr[j-1].idx)
			}
		}
	}
	mpY := make(map[int][]pair)
	for i, m := range mines {
		mpY[m.y] = append(mpY[m.y], pair{m.x, i})
	}
	for _, arr := range mpY {
		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		for j := 1; j < len(arr); j++ {
			if arr[j].val-arr[j-1].val <= k {
				d.union(arr[j].idx, arr[j-1].idx)
			}
		}
	}
	compMin := make(map[int]int)
	maxTime := 0
	for i, m := range mines {
		root := d.find(i)
		if val, ok := compMin[root]; !ok || m.t < val {
			compMin[root] = m.t
		}
		if m.t > maxTime {
			maxTime = m.t
		}
	}
	times := make([]int, 0, len(compMin))
	for _, t := range compMin {
		times = append(times, t)
	}
	sort.Ints(times)
	hi := maxTime
	if hi < len(times) {
		hi = len(times)
	}
	l, r := -1, hi
	for r-l > 1 {
		mid := (l + r) / 2
		idx := sort.Search(len(times), func(i int) bool { return times[i] > mid })
		cnt := len(times) - idx
		if cnt <= mid+1 {
			r = mid
		} else {
			l = mid
		}
	}
	return r
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	k := rng.Intn(5)
	mines := make([]mine, n)
	for i := 0; i < n; i++ {
		mines[i] = mine{
			x: rng.Intn(21) - 10,
			y: rng.Intn(21) - 10,
			t: rng.Intn(20),
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
	for _, m := range mines {
		fmt.Fprintf(&sb, "%d %d %d\n", m.x, m.y, m.t)
	}
	return sb.String(), fmt.Sprint(solveCase(mines, k))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
