package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a [14]int64
	for i := 0; i < 14; i++ {
		fmt.Fscan(in, &a[i])
	}
	var best int64
	for i := 0; i < 14; i++ {
		if a[i] == 0 {
			continue
		}
		b := make([]int64, 14)
		copy(b, a[:])
		x := b[i]
		b[i] = 0
		add := x / 14
		for j := 0; j < 14; j++ {
			b[j] += add
		}
		r := int(x % 14)
		for j := 1; j <= r; j++ {
			idx := (i + j) % 14
			b[idx]++
		}
		var cur int64
		for j := 0; j < 14; j++ {
			if b[j]%2 == 0 {
				cur += b[j]
			}
		}
		if cur > best {
			best = cur
		}
	}
	fmt.Println(best)
}
