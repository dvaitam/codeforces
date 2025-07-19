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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for t := 0; t < T; t++ {
		var n, l int
		fmt.Fscan(reader, &n, &l)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		ans := 0
		for i := 0; i < 30; i++ {
			cnt := 0
			for j := 0; j < n; j++ {
				if (a[j]>>i)&1 == 1 {
					cnt++
				}
			}
			if cnt*2 >= n {
				ans |= 1 << i
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
