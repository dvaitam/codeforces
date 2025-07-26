package main

import (
	"bufio"
	"fmt"
	"os"
)

type fenwick struct {
	f []uint8
	n int
}

func newFenwick(n int) *fenwick {
	return &fenwick{make([]uint8, n+2), n}
}

func (fw *fenwick) add(idx int) {
	idx++
	for idx <= fw.n+1 {
		fw.f[idx] ^= 1
		idx += idx & -idx
	}
}

func (fw *fenwick) sum(idx int) uint8 {
	idx++
	var res uint8
	for idx > 0 {
		res ^= fw.f[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			grid[i] = []byte(s)
		}
		offset := n - 1
		size := 2*n - 1
		fw := newFenwick(size)
		ans := 0
		for s := 2; s <= 2*n; s++ {
			iLow := 1
			if s-n > iLow {
				iLow = s - n
			}
			iHigh := n
			if s-1 < iHigh {
				iHigh = s - 1
			}
			for i := iLow; i <= iHigh; i++ {
				j := s - i
				idx := (i - j) + offset
				val := grid[i-1][j-1] - '0'
				if fw.sum(idx)%2 == 1 {
					val ^= 1
				}
				if val == 1 {
					ans++
					fw.add(idx)
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
