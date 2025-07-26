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
		var rating int
		fmt.Fscan(in, &rating)
		switch {
		case rating >= 1900:
			fmt.Fprintln(out, "Division 1")
		case rating >= 1600:
			fmt.Fprintln(out, "Division 2")
		case rating >= 1400:
			fmt.Fprintln(out, "Division 3")
		default:
			fmt.Fprintln(out, "Division 4")
		}
	}
}
