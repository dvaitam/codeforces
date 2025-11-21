package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd64(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func missingCount(x, g int64, n int64) int64 {
	q := x / g
	res := x - q
	extra := q + 1 - n
	if extra > 0 {
		res += extra
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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		if n == 1 {
			val := a[0]
			k64 := int64(k)
			var ans int64
			if k64 <= val {
				ans = k64 - 1
			} else {
				ans = k64
			}
			fmt.Fprintln(writer, ans)
			continue
		}

		g := a[0]
		for i := 1; i < n; i++ {
			g = gcd64(g, a[i])
		}

		n64 := int64(n)
		k64 := int64(k)

		lo := int64(0)
		hi := g * (n64 + k64 + 5)
		if hi < 0 {
			hi = (n64 + k64 + 5)
		}
		for lo < hi {
			mid := (lo + hi) / 2
			if missingCount(mid, g, n64) >= k64 {
				hi = mid
			} else {
				lo = mid + 1
			}
		}
		fmt.Fprintln(writer, lo)
	}
}
