package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	words := (m + 63) >> 6
	grid := make([][]uint64, n)
	for i := range grid {
		grid[i] = make([]uint64, words)
	}

	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		col := 0
		for j := 0; j < len(line); j++ {
			ch := line[j]
			var v int
			if ch >= '0' && ch <= '9' {
				v = int(ch - '0')
			} else {
				v = int(ch-'A') + 10
			}
			for b := 3; b >= 0; b-- {
				if (v>>uint(b))&1 == 1 {
					w := col >> 6
					bit := uint(col & 63)
					grid[i][w] |= 1 << bit
				}
				col++
			}
		}
	}

	visited := make([][]uint64, n)
	for i := range visited {
		visited[i] = make([]uint64, words)
	}

	encode := func(r, c int) uint32 { return uint32((r << 14) | c) }
	mask := uint32((1 << 14) - 1)

	getBit := func(mat [][]uint64, r, c int) bool {
		return (mat[r][c>>6]>>(uint(c)&63))&1 == 1
	}
	setBit := func(mat [][]uint64, r, c int) {
		mat[r][c>>6] |= 1 << (uint(c) & 63)
	}

	queue := make([]uint32, 0)
	components := 0
	for r := 0; r < n; r++ {
		for w := 0; w < words; w++ {
			bitsLeft := grid[r][w] &^ visited[r][w]
			for bitsLeft != 0 {
				lsb := bitsLeft & -bitsLeft
				c := w*64 + bits.TrailingZeros64(lsb)
				visited[r][w] |= lsb
				queue = queue[:0]
				queue = append(queue, encode(r, c))
				for head := 0; head < len(queue); head++ {
					idx := queue[head]
					x := int(idx >> 14)
					y := int(idx & mask)
					if x > 0 && getBit(grid, x-1, y) && !getBit(visited, x-1, y) {
						setBit(visited, x-1, y)
						queue = append(queue, encode(x-1, y))
					}
					if x+1 < n && getBit(grid, x+1, y) && !getBit(visited, x+1, y) {
						setBit(visited, x+1, y)
						queue = append(queue, encode(x+1, y))
					}
					if y > 0 && getBit(grid, x, y-1) && !getBit(visited, x, y-1) {
						setBit(visited, x, y-1)
						queue = append(queue, encode(x, y-1))
					}
					if y+1 < m && getBit(grid, x, y+1) && !getBit(visited, x, y+1) {
						setBit(visited, x, y+1)
						queue = append(queue, encode(x, y+1))
					}
				}
				components++
				bitsLeft &^= lsb
			}
		}
	}

	fmt.Fprintln(out, components)
}
