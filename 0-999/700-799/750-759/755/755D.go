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

	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	if k > n-k {
		k = n - k
	}
	sections := make([]int64, n)
	var f int64 = 1
	for t := int64(1); t <= n; t++ {
		var inter int64
		if t < n {
			inter = (t*k)/n + ((t-1)*k)/n
		} else {
			inter = 2 * (k - 1)
		}
		f += inter + 1
		sections[t-1] = f
	}
	for i, v := range sections {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
