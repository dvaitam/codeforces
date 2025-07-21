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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		cnt := make([]int, 101)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x >= 0 && x < len(cnt) {
				cnt[x]++
			}
		}
		odd := false
		for i := 1; i <= 50; i++ {
			if cnt[i]%2 == 1 {
				odd = true
				break
			}
		}
		if odd {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
