package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solve(k, n1, n2, n3, t1, t2, t3 int) string {
	wheap := &IntHeap{}
	heap.Init(wheap)
	for i := 0; i < n1; i++ {
		heap.Push(wheap, 0)
	}
	wfin := make([]int, k)
	for i := 0; i < k; i++ {
		avail := heap.Pop(wheap).(int)
		finish := avail + t1
		wfin[i] = finish
		heap.Push(wheap, finish)
	}
	dheap := &IntHeap{}
	heap.Init(dheap)
	for i := 0; i < n2; i++ {
		heap.Push(dheap, 0)
	}
	dfin := make([]int, k)
	for i, wf := range wfin {
		avail := heap.Pop(dheap).(int)
		start := wf
		if avail > start {
			start = avail
		}
		finish := start + t2
		dfin[i] = finish
		heap.Push(dheap, finish)
	}
	fheap := &IntHeap{}
	heap.Init(fheap)
	for i := 0; i < n3; i++ {
		heap.Push(fheap, 0)
	}
	ans := 0
	for _, df := range dfin {
		avail := heap.Pop(fheap).(int)
		start := df
		if avail > start {
			start = avail
		}
		finish := start + t3
		if finish > ans {
			ans = finish
		}
		heap.Push(fheap, finish)
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(20) + 1
	n1 := rng.Intn(5) + 1
	n2 := rng.Intn(5) + 1
	n3 := rng.Intn(5) + 1
	t1 := rng.Intn(10) + 1
	t2 := rng.Intn(10) + 1
	t3 := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d %d %d %d %d %d\n", k, n1, n2, n3, t1, t2, t3)
	expected := solve(k, n1, n2, n3, t1, t2, t3)
	return input, expected
}

func runCase(bin, input, expected string) error {
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
