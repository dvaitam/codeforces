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
		var n int
		fmt.Fscan(reader, &n)
		cnt := make([]int, 101)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x >= 0 && x <= 100 {
				cnt[x]++
			}
		}
		first := true
		for i := 0; i <= 100; i++ {
			if cnt[i] > 0 {
				if !first {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, i)
				cnt[i]--
				first = false
			}
		}
		for i := 0; i <= 100; i++ {
			for cnt[i] > 0 {
				if !first {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, i)
				cnt[i]--
				first = false
			}
		}
		fmt.Fprintln(writer)
	}
}
