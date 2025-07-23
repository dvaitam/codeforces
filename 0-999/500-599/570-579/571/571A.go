package main

import (
	"bufio"
	"fmt"
	"os"
)

func choose2(n int64) int64 {
	return n * (n - 1) / 2
}

func choose3(n int64) int64 {
	return n * (n - 1) * (n - 2) / 6
}

func invalid(x, y, z, l int64) int64 {
	// invalid configurations where side x >= y+z
	res := int64(0)
	for add := int64(0); add <= l; add++ {
		excess := x + add - (y + z)
		if excess < 0 {
			continue
		}
		remain := l - add
		if excess < remain {
			remain = excess
		}
		if remain >= 0 {
			res += choose2(remain + 2)
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b, c, l int64
	if _, err := fmt.Fscan(in, &a, &b, &c, &l); err != nil {
		return
	}
	total := choose3(l + 3)
	ans := total
	ans -= invalid(a, b, c, l)
	ans -= invalid(b, a, c, l)
	ans -= invalid(c, a, b, l)
	fmt.Println(ans)
}
