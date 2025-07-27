package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const maxV = 1 << 20

var freq [maxV]int
var bitset [maxV >> 6]uint64 // since maxV is power of two

func setBit(v int) {
	bitset[v>>6] |= 1 << (v & 63)
}
func clearBit(v int) {
	bitset[v>>6] &^= 1 << (v & 63)
}
func hasBit(v int) bool {
	return bitset[v>>6]&(1<<(v&63)) != 0
}

func nextSet(start int) int {
	if start >= maxV {
		return -1
	}
	idx := start >> 6
	word := bitset[idx] >> (start & 63)
	if word != 0 {
		return (idx << 6) + (start & 63) + int(bits.TrailingZeros64(word))
	}
	idx++
	for idx < len(bitset) {
		if bitset[idx] != 0 {
			return (idx << 6) + int(bits.TrailingZeros64(bitset[idx]))
		}
		idx++
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		freq[v]++
		setBit(v)
	}

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t >= 1 && t <= 3 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			positions := make([]int, 0)
			for i := nextSet(l); i != -1 && i <= r; i = nextSet(i + 1) {
				positions = append(positions, i)
			}
			for _, v := range positions {
				cnt := freq[v]
				if cnt == 0 {
					continue
				}
				freq[v] = 0
				clearBit(v)
				var nv int
				if t == 1 {
					nv = v & x
				} else if t == 2 {
					nv = v | x
				} else {
					nv = v ^ x
				}
				freq[nv] += cnt
				setBit(nv)
			}
		} else if t == 4 {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			cnt := 0
			for i := nextSet(l); i != -1 && i <= r; i = nextSet(i + 1) {
				cnt++
			}
			fmt.Fprintln(writer, cnt)
		}
	}
}
