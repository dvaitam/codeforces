package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		total := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			total ^= arr[i]
		}
		if total == 0 {
			fmt.Fprintln(out, "YES")
			continue
		}
		prefix := 0
		cnt := 0
		for i := 0; i < n; i++ {
			prefix ^= arr[i]
			if prefix == total {
				cnt++
				prefix = 0
			}
		}
		if cnt >= 2 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
