package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var m, k, a1, ak int
		fmt.Fscan(reader, &m, &k, &a1, &ak)

		q := m / k
		r := m % k

		use1 := a1
		if use1 > r {
			use1 = r
		}
		r -= use1
		a1 -= use1

		if ak >= q {
			fmt.Fprintln(writer, r)
			continue
		}
		q -= ak

		replace := a1 / k
		if replace > q {
			replace = q
		}
		q -= replace

		fmt.Fprintln(writer, q+r)
	}
}
