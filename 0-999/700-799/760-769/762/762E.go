package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// BIT is a Fenwick tree supporting point updates and prefix sums.
type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	if i > b.n {
		i = b.n
	}
	s := 0
	for i > 0 {
		s += b.tree[i]
		i &= i - 1
	}
	return s
}

// Station describes a radio station.
type Station struct {
	x   int
	r   int
	f   int
	idx int // compressed index of x in its frequency list
}

// uniqueInts removes consecutive duplicates from a sorted slice.
func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	stations := make([]Station, n)
	freqCoords := make(map[int][]int)

	for i := 0; i < n; i++ {
		var x, r, f int
		fmt.Fscan(in, &x, &r, &f)
		stations[i] = Station{x: x, r: r, f: f}
		freqCoords[f] = append(freqCoords[f], x)
	}

	// Prepare coordinate compression and BIT for each frequency
	bits := make(map[int]*BIT, len(freqCoords))
	coords := make(map[int][]int, len(freqCoords))
	for f, arr := range freqCoords {
		sort.Ints(arr)
		arr = uniqueInts(arr)
		coords[f] = arr
		bits[f] = NewBIT(len(arr))
	}

	// assign compressed index for each station
	for i := range stations {
		arr := coords[stations[i].f]
		stations[i].idx = sort.SearchInts(arr, stations[i].x) + 1
	}

	// sort stations by x coordinate
	sort.Slice(stations, func(i, j int) bool { return stations[i].x < stations[j].x })

	ans := int64(0)
	left := 0 // pointer to the first active station
	for i, st := range stations {
		// remove stations whose coverage does not reach st.x anymore
		for left < i && stations[left].x+stations[left].r < st.x {
			rm := stations[left]
			bits[rm.f].Add(rm.idx, -1)
			left++
		}

		// count bad pairs with station st as the right endpoint
		for freq := st.f - k; freq <= st.f+k; freq++ {
			if freq < 1 || freq > 10000 {
				continue
			}
			bit := bits[freq]
			if bit == nil {
				continue
			}
			arr := coords[freq]
			// find first index with coordinate >= st.x - st.r
			pos := sort.SearchInts(arr, st.x-st.r)
			cnt := bit.Sum(len(arr)) - bit.Sum(pos)
			ans += int64(cnt)
		}

		// insert current station into its frequency structure
		bits[st.f].Add(st.idx, 1)
	}

	fmt.Fprintln(out, ans)
}
