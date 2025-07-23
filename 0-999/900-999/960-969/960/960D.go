package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func mod(x, m int64) int64 {
	x %= m
	if x < 0 {
		x += m
	}
	return x
}

var shiftVal [61]int64
var shiftNode [61]int64

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var x, k int64
			fmt.Fscan(reader, &x, &k)
			l := bits.Len64(uint64(x)) - 1
			m := int64(1) << uint(l)
			k = mod(k, m)
			shiftVal[l] = mod(shiftVal[l]+k, m)
		} else if t == 2 {
			var x, k int64
			fmt.Fscan(reader, &x, &k)
			l := bits.Len64(uint64(x)) - 1
			m := int64(1) << uint(l)
			k = mod(k, m)
			shiftNode[l] = mod(shiftNode[l]+k, m)
		} else {
			var x int64
			fmt.Fscan(reader, &x)
			l := bits.Len64(uint64(x)) - 1
			m := int64(1) << uint(l)
			idx := mod(x-m+shiftVal[l]+shiftNode[l], m)
			first := true
			for {
				val := (int64(1) << uint(l)) + mod(idx-shiftNode[l]-shiftVal[l], m)
				if !first {
					writer.WriteByte(' ')
				}
				first = false
				fmt.Fprint(writer, val)
				if l == 0 {
					break
				}
				idx >>= 1
				l--
				m = int64(1) << uint(l)
				idx = mod(idx+shiftNode[l], m)
			}
			writer.WriteByte('\n')
		}
	}
}
