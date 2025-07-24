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
		var n, m int64
		fmt.Fscan(reader, &n, &m)

		k := int64(1)
		cnt2, cnt5 := 0, 0
		tmp := n
		for tmp%2 == 0 {
			cnt2++
			tmp /= 2
		}
		for tmp%5 == 0 {
			cnt5++
			tmp /= 5
		}

		for k*2 <= m && cnt2 < cnt5 {
			k *= 2
			cnt2++
		}
		for k*5 <= m && cnt5 < cnt2 {
			k *= 5
			cnt5++
		}
		for k*10 <= m {
			k *= 10
		}
		k *= m / k
		fmt.Fprintln(writer, n*k)
	}
}
