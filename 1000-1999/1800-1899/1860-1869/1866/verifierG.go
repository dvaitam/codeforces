package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

type segment struct {
	L, R int
	A    int64
}

type item struct {
	R      int
	remain int64
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].R < h[j].R }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, tc := range tests {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		if err := check(tc.expect, strings.TrimSpace(out)); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", idx+1, err, tc.input, tc.expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func check(expect, actual string) error {
	exp, _ := strconv.ParseInt(expect, 10, 64)
	val, err := strconv.ParseInt(actual, 10, 64)
	if err != nil {
		return fmt.Errorf("output is not integer: %v", err)
	}
	if val != exp {
		return fmt.Errorf("expected %d but got %d", exp, val)
	}
	return nil
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeCase([]int64{7, 1, 2}, []int{1, 2, 0}),
		makeCase([]int64{0, 0, 10}, []int{0, 0, 2}),
		makeCase([]int64{5}, []int{0}),
	}
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 1
		A := make([]int64, n)
		D := make([]int, n)
		for j := 0; j < n; j++ {
			A[j] = rand.Int63n(10)
			D[j] = rand.Intn(n)
		}
		tests = append(tests, makeCase(A, D))
	}
	return tests
}

func makeCase(A []int64, D []int) testCase {
	var sb strings.Builder
	n := len(A)
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range D {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return testCase{
		input:  sb.String(),
		expect: fmt.Sprintf("%d", solveReference(A, D)),
	}
}

func solveReference(A []int64, D []int) int64 {
	n := len(A)
	segs := make([]segment, n)
	var total int64
	for i := 0; i < n; i++ {
		L := i + 1 - D[i]
		if L < 1 {
			L = 1
		}
		R := i + 1 + D[i]
		if R > n {
			R = n
		}
		segs[i] = segment{L: L, R: R, A: A[i]}
		total += A[i]
	}
	sort.Slice(segs, func(i, j int) bool { return segs[i].L < segs[j].L })
	low, high := int64(0), total
	for low < high {
		mid := (low + high) / 2
		if feasible(mid, segs, n) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func feasible(z int64, segs []segment, n int) bool {
	idx := 0
	pq := &minHeap{}
	for pos := 1; pos <= n; pos++ {
		for idx < len(segs) && segs[idx].L <= pos {
			if segs[idx].A > 0 {
				*pq = append(*pq, item{R: segs[idx].R, remain: segs[idx].A})
				upHeap(pq, pq.Len()-1)
			}
			idx++
		}
		for pq.Len() > 0 && (*pq)[0].R < pos {
			return false
		}
		capacity := z
		for capacity > 0 && pq.Len() > 0 {
			top := &(*pq)[0]
			if top.R < pos {
				return false
			}
			if top.remain <= capacity {
				capacity -= top.remain
				popHeap(pq)
			} else {
				top.remain -= capacity
				capacity = 0
			}
		}
	}
	return pq.Len() == 0
}

func (h minHeap) Len() int {
	return len(h)
}

func upHeap(h *minHeap, i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.Less(i, parent) {
			break
		}
		h.Swap(i, parent)
		i = parent
	}
}

func popHeap(h *minHeap) {
	last := len(*h) - 1
	h.Swap(0, last)
	*h = (*h)[:last]
	downHeap(h, 0)
}

func downHeap(h *minHeap, i int) {
	n := len(*h)
	for {
		l := 2*i + 1
		r := 2*i + 2
		smallest := i
		if l < n && h.Less(l, smallest) {
			smallest = l
		}
		if r < n && h.Less(r, smallest) {
			smallest = r
		}
		if smallest == i {
			break
		}
		h.Swap(i, smallest)
		i = smallest
	}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}
