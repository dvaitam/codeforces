package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const (
	limitNode = 450
	maxDelta  = 305000
)

var (
	anchors  = []int{2, 3, 4}
	nodes    []int
	bitsets  [][]uint64
	words    int
	lastMask uint64
	manual   = map[int][][2]int{
		3: {{1, 3}, {2, 3}},
		4: {{1, 2}, {1, 3}, {1, 4}},
		5: {{3, 1}, {3, 2}, {3, 4}, {3, 5}},
		6: {{3, 1}, {4, 1}, {5, 2}, {2, 1}, {1, 6}},
		7: {{2, 1}, {3, 1}, {6, 1}, {1, 5}, {5, 4}, {4, 7}},
	}
)

func init() {
	for i := 5; i <= limitNode; i++ {
		nodes = append(nodes, i)
	}
	words = (maxDelta + 64) / 64
	rem := (maxDelta + 1) % 64
	if rem == 0 {
		lastMask = ^uint64(0)
	} else {
		lastMask = (uint64(1) << uint(rem)) - 1
	}
	buildBitsets()
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n == 2 {
			fmt.Fprintln(out, -1)
			continue
		}
		if n <= 7 {
			edges := manual[n]
			for _, e := range edges {
				fmt.Fprintf(out, "%d %d\n", e[0], e[1])
			}
			continue
		}

		parents := make([]int, n+1)
		for i := 2; i <= n; i++ {
			parents[i] = 1
		}

		base := int64(n)*(int64(n)+1)/2 - 1
		r := ceilSqrt(base)
		prefix := min(n, limitNode) - 4
		if prefix < 0 {
			prefix = 0
		}
		delta := 0
		for tries := 0; tries < 1000; tries++ {
			delta = int(r*r - base)
			if delta == 0 {
				break
			}
			if delta == 1 || delta == 2 || delta == 4 {
				r++
				continue
			}
			if delta > maxDelta {
				r++
				continue
			}
			if bitTest(bitsets[prefix], delta) {
				break
			}
			r++
		}
		if delta > 0 {
			applyAssignments(prefix, delta, n, parents)
		}
		for i := 2; i <= n; i++ {
			fmt.Fprintf(out, "%d %d\n", i, parents[i])
		}
	}
}

func ceilSqrt(x int64) int64 {
	r := int64(math.Sqrt(float64(x)))
	for r*r < x {
		r++
	}
	for (r-1)*(r-1) >= x {
		r--
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildBitsets() {
	bitsets = make([][]uint64, len(nodes)+1)
	cur := make([]uint64, words)
	cur[0] = 1
	bitsets[0] = copySlice(cur)

	base := make([]uint64, words)
	tmp := make([]uint64, words)

	for idx, node := range nodes {
		copy(base, cur)
		for _, anchor := range anchors {
			if anchor == node {
				continue
			}
			shift := (anchor - 1) * node
			if shift > maxDelta {
				continue
			}
			shiftLeftInto(tmp, base, shift)
			for i := 0; i < words; i++ {
				cur[i] |= tmp[i]
			}
		}
		cur[words-1] &= lastMask
		bitsets[idx+1] = copySlice(cur)
	}
}

func shiftLeftInto(dst, src []uint64, shift int) {
	for i := range dst {
		dst[i] = 0
	}
	wordShift := shift / 64
	bitShift := uint(shift % 64)
	for i := len(src) - 1; i >= 0; i-- {
		j := i - wordShift
		if j < 0 {
			continue
		}
		val := src[j] << bitShift
		if bitShift != 0 && j > 0 {
			val |= src[j-1] >> (64 - bitShift)
		}
		dst[i] = val
	}
}

func copySlice(src []uint64) []uint64 {
	dup := make([]uint64, len(src))
	copy(dup, src)
	return dup
}

func bitTest(bits []uint64, pos int) bool {
	if pos < 0 {
		return false
	}
	idx := pos / 64
	if idx >= len(bits) {
		return false
	}
	off := uint(pos % 64)
	return (bits[idx]>>off)&1 == 1
}

func applyAssignments(prefix, delta, n int, parents []int) {
	remain := delta
	for idx := prefix; idx >= 1; idx-- {
		node := nodes[idx-1]
		if node > n {
			continue
		}
		prevBits := bitsets[idx-1]
		if bitTest(prevBits, remain) {
			continue
		}
		assigned := false
		for _, anchor := range anchors {
			if anchor == node {
				continue
			}
			val := node * (anchor - 1)
			if val > remain {
				continue
			}
			if bitTest(prevBits, remain-val) {
				parents[node] = anchor
				remain -= val
				assigned = true
				break
			}
		}
		if !assigned {
			panic("assignment failed")
		}
	}
	if remain != 0 {
		panic("incomplete assignment")
	}
}
