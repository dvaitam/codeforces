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
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	var g int64
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		if i == 0 {
			g = x
		} else {
			g = gcd(g, x)
		}
	}

	ans := int64(1)
	for p := int64(2); p*p <= g; p++ {
		if g%p == 0 {
			cnt := int64(0)
			for g%p == 0 {
				g /= p
				cnt++
			}
			ans *= cnt + 1
		}
	}
	if g > 1 {
		ans *= 2
	}

	fmt.Fprintln(writer, ans)
}
