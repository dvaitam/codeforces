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

	var n, k, m, d int64
	if _, err := fmt.Fscan(reader, &n, &k, &m, &d); err != nil {
		return
	}
	var ans int64
	for i := int64(1); i <= d; i++ {
		t := (i-1)*k + 1
		if t > n {
			break
		}
		x := n / t
		if x > m {
			x = m
		}
		if x == 0 {
			continue
		}
		val := x * i
		if val > ans {
			ans = val
		}
	}
	fmt.Fprintln(writer, ans)
}
