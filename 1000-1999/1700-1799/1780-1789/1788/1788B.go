package main

import (
	"bufio"
	"fmt"
	"os"
)

// Precomputed adjustments for trailing 9 counts
var ans = [10]int64{0, 1, 2, 10, 18, 100, 180, 1000, 1800, 10000}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for i := 0; i < t; i++ {
		var n int64
		fmt.Fscan(reader, &n)
		var a, b int64
		if n%2 == 0 {
			a = n / 2
			b = n / 2
		} else if n%10 != 9 {
			a = n / 2
			b = n - a
		} else {
			p := n / 2
			q := n - p
			tmp := n
			cnt := 0
			for tmp%10 == 9 {
				tmp /= 10
				cnt++
			}
			val := tmp % 10
			if val == 0 {
				val = 9
				cnt--
			}
			var offset int64
			if val%2 == 0 {
				offset = ans[cnt-1] * 5
			} else {
				offset = ans[cnt] * 5
			}
			a = p + offset
			b = q - offset
		}
		fmt.Fprintln(writer, a, b)
	}
}
