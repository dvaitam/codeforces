package main

import (
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

func run(bin, input string) (string, error) {
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

type task struct {
	t, s, p int
}

type info struct {
	x, y int
}

type maxHeap []info

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i].x > h[j].x }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(info)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func calcOverlap(s1, t1, s2, t2 int64) int64 {
	if s2 > s1 {
		s1 = s2
	}
	if t2 < t1 {
		t1 = t2
	}
	res := t1 - s1 + 1
	if res > 0 {
		return res
	}
	return 0
}

func simulate(tasks []task, order []info, w int, T int64, updateUsed bool, used, last []int64) {
	var cur int64
	h := &maxHeap{}
	heap.Init(h)
	need := make([]int, len(tasks))
	for i := 0; i < len(tasks); i++ {
		heap.Push(h, info{tasks[order[i].y].p, order[i].y})
		need[order[i].y] = tasks[order[i].y].s
		if cur < int64(order[i].x) {
			cur = int64(order[i].x)
		}
		for (i == len(tasks)-1 || cur < int64(order[i+1].x)) && h.Len() > 0 {
			top := (*h)[0]
			x := top.y
			var d int
			if i == len(tasks)-1 {
				d = need[x]
			} else {
				avail := int(int64(order[i+1].x) - cur)
				if need[x] < avail {
					d = need[x]
				} else {
					d = avail
				}
			}
			added := calcOverlap(int64(tasks[w].t), T-1, cur, cur+int64(d)-1)
			if updateUsed {
				used[x] += added
			}
			cur += int64(d)
			need[x] -= d
			if need[x] == 0 {
				last[x] = cur
				heap.Pop(h)
			}
		}
	}
}

func solveE(n int, tasks []task, T int64, w int) (int, []int64) {
	order := make([]info, n)
	for i := 0; i < n; i++ {
		order[i] = info{tasks[i].t, i}
	}
	sort.Slice(order, func(i, j int) bool { return order[i].x < order[j].x })
	used := make([]int64, n)
	last := make([]int64, n)
	simulate(tasks, order, w, T, true, used, last)
	g := make([]info, n)
	for i := 0; i < n; i++ {
		g[i] = info{tasks[i].p, i}
	}
	sort.Slice(g, func(i, j int) bool { return g[i].x < g[j].x })
	var sum int64
	j := 0
	for sum < int64(tasks[w].s) {
		sum += used[g[j].y]
		j++
	}
	j--
	newP := g[j].x + 1
	if newP < 1 {
		newP = 1
	}
	for j+1 < n && g[j+1].x == newP {
		newP++
		j++
	}
	tasks[w].p = newP
	simulate(tasks, order, w, T, false, used, last)
	return newP, last
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	tasks := make([]task, n)
	usedPriorities := map[int]bool{}
	for i := 0; i < n; i++ {
		tasks[i].t = rng.Intn(20)
		tasks[i].s = rng.Intn(5) + 1
		p := rng.Intn(30) + 1
		for usedPriorities[p] {
			p = rng.Intn(30) + 1
		}
		tasks[i].p = p
		usedPriorities[p] = true
	}
	w := rng.Intn(n)
	order := make([]info, n)
	for i := 0; i < n; i++ {
		order[i] = info{tasks[i].t, i}
	}
	sort.Slice(order, func(i, j int) bool { return order[i].x < order[j].x })
	used := make([]int64, n)
	last := make([]int64, n)
	simulate(tasks, order, w, 1<<60, false, used, last)
	T := last[w]
	tasks[w].p = -1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tasks[i].t, tasks[i].s, tasks[i].p))
	}
	sb.WriteString(fmt.Sprintf("%d\n", T))
	input := sb.String()
	expectP, expectLast := solveE(n, append([]task(nil), tasks...), T, w)
	expectStr := fmt.Sprintf("%d\n", expectP)
	for i := 0; i < n; i++ {
		expectStr += fmt.Sprintf("%d ", expectLast[i])
	}
	expectStr = strings.TrimSpace(expectStr)
	return input, expectStr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tc, exp, strings.TrimSpace(out), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
