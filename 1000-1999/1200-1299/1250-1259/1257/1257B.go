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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x >= y {
			fmt.Fprintln(out, "YES")
			continue
		}
		if x == 1 {
			fmt.Fprintln(out, "NO")
		} else if x == 2 && y == 3 {
			fmt.Fprintln(out, "YES")
		} else if x <= 3 {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}
