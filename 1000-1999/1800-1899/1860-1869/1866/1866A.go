package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ans := int(^uint(0) >> 1) // max int
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		v := abs(a)
		if v < ans {
			ans = v
		}
	}
	fmt.Fprintln(writer, ans)
}
