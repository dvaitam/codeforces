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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	var sum int
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		sum ^= a[i]
	}
	if n%2 == 1 {
		fmt.Fprintln(out, "YES")
		fmt.Fprintln(out, n-2)
		// forward triples
		for i := 0; i+2 < n; i += 2 {
			fmt.Fprintf(out, "%d %d %d\n", i+1, i+2, i+3)
		}
		// backward triples
		for i := n - 3; i-2 >= 0; i -= 2 {
			fmt.Fprintf(out, "%d %d %d\n", i+1, (i-1)+1, (i-2)+1)
		}
	} else {
		if sum == 0 {
			fmt.Fprintln(out, "YES")
			fmt.Fprintln(out, n-3)
			// exclude last element
			n--
			for i := 0; i+2 < n; i += 2 {
				fmt.Fprintf(out, "%d %d %d\n", i+1, i+2, i+3)
			}
			for i := n - 3; i-2 >= 0; i -= 2 {
				fmt.Fprintf(out, "%d %d %d\n", i+1, (i-1)+1, (i-2)+1)
			}
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
