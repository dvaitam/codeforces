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

type Node struct {
	now  int64
	x, y int64
	shu  int64
	flag bool
}

type MaxHeap []Node

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].now > h[j].now }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Node)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func runCandidate(bin, input string) (string, error) {
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

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func get(n, m, r, x, y int64) int64 {
	lx := max(1, x-r+1)
	ly := max(1, y-r+1)
	x = min(x, n-r+1)
	y = min(y, m-r+1)
	return (x - lx + 1) * (y - ly + 1)
}

func solveCase(n, m, r, k int64) string {
	number := (n - r + 1) * (m - r + 1)
	MAX := get(n, m, r, (n+1)/2, (m+1)/2)
	var a, b int64
	if n >= 2*r {
		a = n - 2*r + 2
	} else {
		a = 2*r - n
	}
	if m >= 2*r {
		b = m - 2*r + 2
	} else {
		b = 2*r - m
	}
	hen := min(n-r+1, r)
	shu := min(m-r+1, r)

	h := &MaxHeap{}
	heap.Init(h)
	heap.Push(h, Node{now: MAX, x: a, y: b, shu: shu, flag: true})

	var ans int64
	kk := k
	for {
		p := heap.Pop(h).(Node)
		cnt := p.x * p.y
		if cnt >= kk {
			ans += p.now * kk
			break
		}
		kk -= cnt
		ans += p.now * cnt
		if p.flag && p.now-hen >= 1 {
			heap.Push(h, Node{now: p.now - hen, x: p.x, y: 2, shu: p.shu - 1, flag: true})
		}
		if p.now-p.shu >= 1 {
			heap.Push(h, Node{now: p.now - p.shu, x: 2, y: p.y, shu: p.shu, flag: false})
		}
	}
	return fmt.Sprintf("%.15f", float64(ans)/float64(number))
}

func generateCase(rng *rand.Rand) (int64, int64, int64, int64) {
	n := int64(rng.Intn(10) + 1)
	m := int64(rng.Intn(10) + 1)
	r := int64(rng.Intn(int(min(n, m))) + 1)
	k := int64(rng.Intn(int(min(n*m, 10))) + 1)
	return n, m, r, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, r, k := generateCase(rng)
		input := fmt.Sprintf("%d %d %d %d\n", n, m, r, k)
		expect := solveCase(n, m, r, k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
