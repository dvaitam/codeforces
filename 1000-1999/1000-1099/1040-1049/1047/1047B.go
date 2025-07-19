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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var x, y int
	ans := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x, &y)
		if x+y > ans {
			ans = x + y
		}
	}
	fmt.Fprint(writer, ans)
}
