package main

import (
	"bufio"
	"fmt"
	"os"
)

var aMap = []int64{0, 2, 3, 1}
var bMap = []int64{0, 3, 1, 2}

func tripleByIndex(idx int64) (int64, int64, int64) {
	k := 0
	size := int64(1)
	tmp := idx
	for tmp >= size {
		tmp -= size
		size *= 4
		k++
	}
	base := int64(1) << (2 * uint(k))
	a := base + tmp

	bOff := int64(0)
	cOff := int64(0)
	x := tmp
	for i := 0; i < k; i++ {
		d := x % 4
		x /= 4
		shift := int64(1) << (2 * uint(i))
		bOff += aMap[d] * shift
		cOff += bMap[d] * shift
	}

	b := 2*base + bOff
	c := 3*base + cOff
	return a, b, c
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		idx := (n - 1) / 3
		pos := (n - 1) % 3
		a, b, c := tripleByIndex(idx)
		var ans int64
		if pos == 0 {
			ans = a
		} else if pos == 1 {
			ans = b
		} else {
			ans = c
		}
		fmt.Fprintln(writer, ans)
	}
}
