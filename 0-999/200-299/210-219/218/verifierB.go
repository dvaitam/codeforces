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

type testCase struct {
	n, m        int
	seats       []int
	input       string
	expectedMax int
	expectedMin int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(1000) + 1
	m := rng.Intn(10) + 1
	seats := make([]int, m)
	total := 0
	for i := 0; i < m; i++ {
		seats[i] = rng.Intn(1000) + 1
		total += seats[i]
	}
	if total < n {
		seats[0] += n - total
		total = n
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", seats[i]))
	}
	sb.WriteByte('\n')
	maxSum, minSum := calc(seats, n)
	return testCase{n: n, m: m, seats: seats, input: sb.String(), expectedMax: maxSum, expectedMin: minSum}
}

func calc(seats []int, n int) (int, int) {
	maxH := &maxHeap{}
	minH := &minHeap{}
	for _, v := range seats {
		heap.Push(maxH, v)
		heap.Push(minH, v)
	}
	heap.Init(maxH)
	heap.Init(minH)
	maxSum, minSum := 0, 0
	for i := 0; i < n; i++ {
		x := heap.Pop(maxH).(int)
		maxSum += x
		if x-1 > 0 {
			heap.Push(maxH, x-1)
		}
		y := heap.Pop(minH).(int)
		minSum += y
		if y-1 > 0 {
			heap.Push(minH, y-1)
		}
	}
	return maxSum, minSum
}

type maxHeap []int

func (h maxHeap) Len() int            { return len(h) }
func (h maxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h maxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type minHeap []int

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func runCase(exe string, tc testCase) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 numbers got %d", len(fields))
	}
	var a, b int
	if _, err := fmt.Sscan(fields[0], &a); err != nil {
		return fmt.Errorf("bad number %q", fields[0])
	}
	if _, err := fmt.Sscan(fields[1], &b); err != nil {
		return fmt.Errorf("bad number %q", fields[1])
	}
	if a != tc.expectedMax || b != tc.expectedMin {
		return fmt.Errorf("expected %d %d got %d %d", tc.expectedMax, tc.expectedMin, a, b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
