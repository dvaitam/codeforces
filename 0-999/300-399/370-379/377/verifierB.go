package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseB struct {
	n, m  int
	s     int64
	bugs  []int
	power []int
	cost  []int
}

func generateCase(rng *rand.Rand) (string, testCaseB) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	bugs := make([]int, m)
	for i := range bugs {
		bugs[i] = rng.Intn(20) + 1
	}
	power := make([]int, n)
	for i := range power {
		power[i] = rng.Intn(20) + 1
	}
	cost := make([]int, n)
	for i := range cost {
		cost[i] = rng.Intn(20)
	}
	s := rng.Int63n(50)
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, s)
	for i, v := range bugs {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	for i, v := range power {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	for i, v := range cost {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	tc := testCaseB{n: n, m: m, s: s, bugs: bugs, power: power, cost: cost}
	return b.String(), tc
}

type req struct{ val, idx int }
type shooter struct{ power, cost, idx int }

type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{}   { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

func minimalDays(tc testCaseB) (int, bool) {
	reqs := make([]req, tc.m)
	for i := 0; i < tc.m; i++ {
		reqs[i] = req{val: tc.bugs[i], idx: i}
	}
	shooters := make([]shooter, tc.n)
	for i := 0; i < tc.n; i++ {
		shooters[i] = shooter{power: tc.power[i], cost: tc.cost[i], idx: i + 1}
	}
	sort.Slice(reqs, func(i, j int) bool { return reqs[i].val < reqs[j].val })
	sort.Slice(shooters, func(i, j int) bool { return shooters[i].power < shooters[j].power })
	chk := func(k int) bool {
		var cst int64
		h := &intHeap{}
		heap.Init(h)
		sp := tc.n - 1
		for i := tc.m - 1; i >= 0; i -= k {
			for sp >= 0 && shooters[sp].power >= reqs[i].val {
				heap.Push(h, shooters[sp].cost)
				sp--
			}
			if h.Len() == 0 {
				return false
			}
			c := heap.Pop(h).(int)
			cst += int64(c)
			if cst > tc.s {
				return false
			}
		}
		return true
	}
	l, r := 1, tc.m+1
	for l < r {
		mid := (l + r) / 2
		if chk(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	if r > tc.m {
		return -1, false
	}
	return r, true
}

func runCase(bin string, input string, tc testCaseB) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := bufio.NewReader(&out)
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return fmt.Errorf("output parse")
	}
	minDays, possible := minimalDays(tc)
	if first == "NO" {
		if possible {
			return fmt.Errorf("expected YES")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("expected YES/NO")
	}
	assign := make([]int, tc.m)
	for i := 0; i < tc.m; i++ {
		if _, err := fmt.Fscan(reader, &assign[i]); err != nil {
			return fmt.Errorf("not enough numbers")
		}
	}
	if !possible {
		return fmt.Errorf("expected NO")
	}
	counts := make([]int, tc.n)
	used := make([]bool, tc.n)
	var costSum int64
	maxCnt := 0
	for i := 0; i < tc.m; i++ {
		id := assign[i]
		if id < 1 || id > tc.n {
			return fmt.Errorf("invalid student")
		}
		if tc.power[id-1] < tc.bugs[i] {
			return fmt.Errorf("student too weak")
		}
		counts[id-1]++
		if counts[id-1] > maxCnt {
			maxCnt = counts[id-1]
		}
		if !used[id-1] {
			used[id-1] = true
			costSum += int64(tc.cost[id-1])
		}
	}
	if costSum > tc.s {
		return fmt.Errorf("exceeds cost")
	}
	if maxCnt != minDays {
		return fmt.Errorf("not minimal days")
	}
	for _, cnt := range counts {
		if cnt > minDays {
			return fmt.Errorf("too many assignments")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
