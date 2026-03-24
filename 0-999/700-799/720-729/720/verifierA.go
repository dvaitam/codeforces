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

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type testCase struct {
	n, m int
	A, B []int
}

func (tc testCase) input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	sb.WriteString(fmt.Sprintf("%d", len(tc.A)))
	for _, v := range tc.A {
		sb.WriteString(fmt.Sprintf(" %d", v))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d", len(tc.B)))
	for _, v := range tc.B {
		sb.WriteString(fmt.Sprintf(" %d", v))
	}
	sb.WriteString("\n")
	return sb.String()
}

type Seat struct {
	D1, D2 int
}

func solve(n, m int, A, B []int) string {
	k := len(A)
	l := len(B)

	seats := make([]Seat, 0, n*m)
	for x := 1; x <= n; x++ {
		for y := 1; y <= m; y++ {
			d1 := x + y
			d2 := x + m + 1 - y
			seats = append(seats, Seat{D1: d1, D2: d2})
		}
	}

	sort.Slice(seats, func(i, j int) bool {
		return seats[i].D1 < seats[j].D1
	})
	sort.Ints(A)
	sort.Ints(B)

	h := &IntHeap{}
	heap.Init(h)

	seatIdx := 0
	for _, a := range A {
		for seatIdx < len(seats) && seats[seatIdx].D1 <= a {
			heap.Push(h, seats[seatIdx].D2)
			seatIdx++
		}
		if h.Len() == 0 {
			return "NO"
		}
		heap.Pop(h)
	}

	rem := make([]int, 0, n*m-k)
	for _, d2 := range *h {
		rem = append(rem, d2)
	}
	for i := seatIdx; i < len(seats); i++ {
		rem = append(rem, seats[i].D2)
	}

	sort.Ints(rem)
	if len(rem) < l {
		return "NO"
	}
	for i := 0; i < l; i++ {
		if rem[i] > B[i] {
			return "NO"
		}
	}

	return "YES"
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	total := n * m
	k := rng.Intn(total + 1)
	l := total - k
	A := make([]int, k)
	B := make([]int, l)
	for i := 0; i < k; i++ {
		A[i] = rng.Intn(n+m) + 1
	}
	for i := 0; i < l; i++ {
		B[i] = rng.Intn(n+m) + 1
	}
	return testCase{n: n, m: m, A: A, B: B}
}

func deterministicCases() []testCase {
	cases := []testCase{
		{n: 1, m: 1, A: []int{2}, B: []int{}},
		{n: 2, m: 2, A: []int{2, 3}, B: []int{2, 3}},
		{n: 3, m: 3, A: []int{2, 3, 4}, B: []int{2, 3, 4, 5, 6}},
	}
	return cases
}

func runCase(bin string, tc testCase) error {
	input := tc.input()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	result := strings.TrimSpace(out.String())
	expect := solve(tc.n, tc.m, append([]int(nil), tc.A...), append([]int(nil), tc.B...))
	if result != expect {
		return fmt.Errorf("expected %s got %s", expect, result)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := deterministicCases()
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
