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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	res := make([]int, 0, m)
	for i := 1; i <= n; i++ {
		lNonWindow := 2*n + 2*(i-1) + 1
		lWindow := 2*(i-1) + 1
		rNonWindow := 2*n + 2*(i-1) + 2
		rWindow := 2*(i-1) + 2
		if lNonWindow <= m {
			res = append(res, lNonWindow)
		}
		if lWindow <= m {
			res = append(res, lWindow)
		}
		if rNonWindow <= m {
			res = append(res, rNonWindow)
		}
		if rWindow <= m {
			res = append(res, rWindow)
		}
	}
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
}
