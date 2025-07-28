package main

import (
	"bytes"
	"container/heap"
	"context"
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
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func solveB(a []int) int64 {
	h := &IntHeap{}
	heap.Init(h)
	var invest, sumTotal, ans int64
	n := len(a)
	for i := 1; i <= n; i++ {
		for invest < int64(i-1) {
			if h.Len() == 0 {
				return ans
			}
			invest += int64(heap.Pop(h).(int))
		}
		if invest < int64(i-1) {
			return ans
		}
		heap.Push(h, a[i-1])
		sumTotal += int64(a[i-1])
		if cur := sumTotal - invest; cur > ans {
			ans = cur
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(2))
	const tests = 100
	for t := 0; t < tests; t++ {
		n := 1 + r.Intn(50)
		a := make([]int, n)
		for i := range a {
			a[i] = r.Intn(n + 1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d error: %v\n", t+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscanf(out, "%d", &got); err != nil {
			fmt.Printf("test %d invalid output\n", t+1)
			os.Exit(1)
		}
		want := solveB(a)
		if got != want {
			fmt.Printf("test %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
