package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// BIT represents a Fenwick tree for prefix sums of ints.
type BIT struct {
	n int
	t []int
}

// NewBIT returns a BIT for indices 1..n
func NewBIT(n int) *BIT {
	return &BIT{n: n, t: make([]int, n+1)}
}

// Add adds delta at index i (1-based)
func (b *BIT) Add(i, delta int) {
	for i <= b.n {
		b.t[i] += delta
		i += i & -i
	}
}

// Sum returns prefix sum up to index i (1-based)
func (b *BIT) Sum(i int) int {
	if i > b.n {
		i = b.n
	}
	s := 0
	for i > 0 {
		s += b.t[i]
		i -= i & -i
	}
	return s
}

// Query stores an operation
type Query struct {
	a int
	t int
	x int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	qs := make([]Query, n)
	times := make(map[int][]int)
	for i := 0; i < n; i++ {
		var a, t, x int
		fmt.Fscan(reader, &a, &t, &x)
		qs[i] = Query{a, t, x}
		times[x] = append(times[x], t)
	}

	// compress times for each value
	comp := make(map[int][]int, len(times))
	bits := make(map[int]*BIT, len(times))
	for v, arr := range times {
		sort.Ints(arr)
		// unique
		m := 1
		for j := 1; j < len(arr); j++ {
			if arr[j] != arr[m-1] {
				arr[m] = arr[j]
				m++
			}
		}
		arr = arr[:m]
		comp[v] = arr
		bits[v] = NewBIT(len(arr))
	}

	// process queries in input order
	for _, q := range qs {
		arr := comp[q.x]
		idx := sort.SearchInts(arr, q.t) + 1 // 1-based
		bit := bits[q.x]
		switch q.a {
		case 1:
			bit.Add(idx, 1)
		case 2:
			bit.Add(idx, -1)
		case 3:
			ans := bit.Sum(idx)
			fmt.Fprintln(writer, ans)
		}
	}
}
