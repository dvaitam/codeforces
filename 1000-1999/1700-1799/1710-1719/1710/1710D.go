package main

import (
	"bufio"
	"fmt"
	"os"
)

func setBit(bs []uint64, idx int) {
	bs[idx>>6] |= 1 << uint(idx&63)
}

func getBit(bs []uint64, idx int) bool {
	return (bs[idx>>6]>>(uint(idx)&63))&1 != 0
}

func maskRange(start, end int) uint64 {
	if start > end {
		return 0
	}
	return (^uint64(0) << uint(start)) & (^uint64(0) >> uint(63-end))
}

func hasSplit(rowI, colJ []uint64, start, end int) bool {
	if start > end {
		return false
	}
	si := start >> 6
	ei := end >> 6
	if si == ei {
		mask := maskRange(start&63, end&63)
		return (rowI[si] & colJ[si] & mask) != 0
	}
	mask := maskRange(start&63, 63)
	if rowI[si]&colJ[si]&mask != 0 {
		return true
	}
	for b := si + 1; b < ei; b++ {
		if rowI[b]&colJ[b] != 0 {
			return true
		}
	}
	mask = maskRange(0, end&63)
	if rowI[ei]&colJ[ei]&mask != 0 {
		return true
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		size := (n + 63) >> 6
		row := make([][]uint64, n+1)
		col := make([][]uint64, n+1)
		for i := 1; i <= n; i++ {
			row[i] = make([]uint64, size)
			col[i] = make([]uint64, size)
		}
		for i := 1; i <= n; i++ {
			var s string
			fmt.Fscan(in, &s)
			for d := 1; d < len(s); d++ {
				if s[d] == '1' {
					j := i + d
					setBit(row[i], j)
					setBit(col[j], i)
				}
			}
		}
		edges := make([][2]int, 0, n-1)
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if !getBit(row[i], j) {
					continue
				}
				if !hasSplit(row[i], col[j], i+1, j-1) {
					edges = append(edges, [2]int{i, j})
				}
			}
		}
		for _, e := range edges {
			fmt.Printf("%d %d\n", e[0], e[1])
		}
	}
}
