package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		ans := 1
		var g int64
		var l int64
		g = 0
		l = 1
		for i := 0; i < n; i++ {
			var a, b int64
			fmt.Fscan(in, &a, &b)
			if g == 0 {
				g = a * b
				l = b
				continue
			}
			newG := gcd(g, a*b)
			newL := l / gcd(l, b) * b
			if newL > newG || newG%newL != 0 {
				ans++
				g = a * b
				l = b
			} else {
				g = newG
				l = newL
			}
		}
		fmt.Fprintln(out, ans)
	}
}
