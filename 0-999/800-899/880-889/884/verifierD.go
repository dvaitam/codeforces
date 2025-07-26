package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type IntHeap []int64

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solve(a []int64) int64 {
	h := &IntHeap{}
	heap.Init(h)
	for _, v := range a {
		heap.Push(h, v)
	}
	n := len(a)
	if n == 1 {
		return 0
	}
	if n%2 == 0 {
		heap.Push(h, 0)
	}
	var cost int64
	for h.Len() > 1 {
		x := heap.Pop(h).(int64)
		y := heap.Pop(h).(int64)
		z := heap.Pop(h).(int64)
		sum := x + y + z
		cost += sum
		heap.Push(h, sum)
	}
	return cost
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(50) + 1
		a := make([]int64, n)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			a[i] = int64(rand.Intn(1000) + 1)
			fmt.Fprintf(&input, "%d ", a[i])
		}
		fmt.Fprintln(&input)
		expected := solve(a)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(expected) {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
