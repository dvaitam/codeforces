package main

import (
	"container/heap"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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

type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func solveCase(arr []int) []int {
	n := len(arr)
	freq := make([]int, n+1)
	for _, v := range arr {
		if v <= n {
			freq[v]++
		}
	}
	fill := make([]int, n)
	for i := range fill {
		fill[i] = -1
	}
	h := &MaxHeap{}
	heap.Init(h)
	cost := 0
	for i := 0; i < n; i++ {
		for c := 0; c < freq[i]; c++ {
			heap.Push(h, i)
		}
		if h.Len() == 0 {
			break
		}
		val := heap.Pop(h).(int)
		cost += i - val
		fill[i] = cost
	}
	ans := make([]int, n+1)
	ans[0] = freq[0]
	impossible := false
	for i := 1; i <= n; i++ {
		if impossible || fill[i-1] == -1 {
			ans[i] = -1
			impossible = true
		} else {
			ans[i] = fill[i-1] + freq[i]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n + 1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	res := solveCase(arr)
	var out strings.Builder
	for i := 0; i <= n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", res[i])
	}
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
