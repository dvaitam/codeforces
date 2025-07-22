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

	var x, y int
	if _, err := fmt.Fscan(in, &x, &y); err != nil {
		return
	}

	if y == 0 {
		fmt.Fprintln(out, "No")
		return
	}

	if y == 1 {
		if x == 0 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
		return
	}

	if x >= y-1 && (x-y+1)%2 == 0 {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
