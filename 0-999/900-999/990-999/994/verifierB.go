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

type knight struct {
	p   int
	c   int64
	idx int
}

type minHeap []int64

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func expected(n, k int, p []int, c []int64) []int64 {
	knights := make([]knight, n)
	for i := 0; i < n; i++ {
		knights[i] = knight{p: p[i], c: c[i], idx: i}
	}
	sort.Slice(knights, func(i, j int) bool { return knights[i].p < knights[j].p })
	ans := make([]int64, n)
	h := &minHeap{}
	heap.Init(h)
	var sum int64
	for _, kn := range knights {
		ans[kn.idx] = sum + kn.c
		heap.Push(h, kn.c)
		sum += kn.c
		if h.Len() > k {
			v := heap.Pop(h).(int64)
			sum -= v
		}
	}
	return ans
}

func runCase(bin string, n, k int, p []int, c []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	expect := expected(n, k, p, c)
	if len(gotFields) != len(expect) {
		return fmt.Errorf("expected %d numbers got %d", len(expect), len(gotFields))
	}
	for i, g := range gotFields {
		var val int64
		if _, err := fmt.Sscan(g, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != expect[i] {
			return fmt.Errorf("expected %v got %v", expect, gotFields)
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int, []int, []int64) {
	n := rng.Intn(50) + 1
	k := rng.Intn(11)
	if k > n-1 {
		k = n - 1
	}
	perm := rng.Perm(n)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = perm[i] + 1 + rng.Intn(100000)
	}
	c := make([]int64, n)
	for i := range c {
		c[i] = int64(rng.Intn(1000))
	}
	return n, k, p, c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic simple cases
	run := []struct {
		n int
		k int
		p []int
		c []int64
	}{
		{1, 0, []int{10}, []int64{5}},
		{3, 1, []int{1, 2, 3}, []int64{10, 20, 30}},
	}

	for i, tc := range run {
		if err := runCase(bin, tc.n, tc.k, tc.p, tc.c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		n, k, p, c := generateCase(rng)
		if err := runCase(bin, n, k, p, c); err != nil {
			fmt.Fprintf(os.Stderr, "random case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
