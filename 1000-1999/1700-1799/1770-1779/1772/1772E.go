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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		asc, desc := 0, 0
		for i := 0; i < n; i++ {
			if arr[i] == i+1 {
				asc++
			}
			if arr[i] == n-i {
				desc++
			}
		}
		if asc > desc {
			fmt.Fprintln(out, "First")
		} else if desc > asc {
			fmt.Fprintln(out, "Second")
		} else {
			fmt.Fprintln(out, "Tie")
		}
	}
}
