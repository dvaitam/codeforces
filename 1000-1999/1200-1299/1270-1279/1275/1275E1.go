package main

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"os"
)

func rawCRC(data []byte) uint32 {
	var crc uint32
	for _, b := range data {
		crc = crc32.IEEETable[byte(crc)^b] ^ (crc >> 8)
	}
	return crc
}

func matrixFor(n, pos int) [32]uint32 {
	var mat [32]uint32
	for bit := 0; bit < 32; bit++ {
		arr := make([]byte, n)
		arr[pos+bit/8] = 1 << uint(bit%8)
		mat[bit] = rawCRC(arr)
	}
	return mat
}

func solveLinear(mat [32]uint32, target uint32) (uint32, bool) {
	var a [32]uint64
	for r := 0; r < 32; r++ {
		row := uint64(0)
		for c := 0; c < 32; c++ {
			if (mat[c]>>uint(r))&1 == 1 {
				row |= 1 << uint(c)
			}
		}
		if (target>>uint(r))&1 == 1 {
			row |= 1 << 32
		}
		a[r] = row
	}
	for c := 0; c < 32; c++ {
		pivot := -1
		for r := c; r < 32; r++ {
			if (a[r]>>uint(c))&1 == 1 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			return 0, false
		}
		if pivot != c {
			a[pivot], a[c] = a[c], a[pivot]
		}
		for r := 0; r < 32; r++ {
			if r != c && ((a[r]>>uint(c))&1 == 1) {
				a[r] ^= a[c]
			}
		}
	}
	var res uint32
	for r := 0; r < 32; r++ {
		if (a[r]>>32)&1 == 1 {
			res |= 1 << uint(r)
		}
	}
	return res, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]byte, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		a[i] = byte(x)
	}

	for ; q > 0; q-- {
		var i, j int
		xs := make([]byte, 4)
		fmt.Fscan(in, &i, &j, &xs[0], &xs[1], &xs[2], &xs[3])

		diffi := make([]byte, n)
		for k := 0; k < 4; k++ {
			diffi[i+k] = a[i+k] ^ xs[k]
		}
		target := rawCRC(diffi)

		mat := matrixFor(n, j)
		sol, ok := solveLinear(mat, target)
		if !ok {
			fmt.Println("No solution")
			continue
		}
		res := make([]byte, 4)
		for k := 0; k < 4; k++ {
			res[k] = a[j+k] ^ byte(sol>>uint(8*k))
		}
		fmt.Printf("%d %d %d %d\n", res[0], res[1], res[2], res[3])
	}
}
