package main

import (
	"bufio"
	"fmt"
	"os"
)

func canFill(n, m int, a []int64) bool {
	total := 0
	hasBig := false
	for _, v := range a {
		cnt := int(v / int64(n))
		if cnt >= 2 {
			total += cnt
			if cnt >= 3 {
				hasBig = true
			}
		}
	}
	if total < m {
		return false
	}
	if m%2 == 1 && !hasBig {
		return false
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		a := make([]int64, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &a[i])
		}
		if canFill(n, m, a) || canFill(m, n, a) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
