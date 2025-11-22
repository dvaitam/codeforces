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
		var start, target int64
		fmt.Fscan(in, &start, &target)
		diff := target - start
		fmt.Fprintf(out, "add %d\n", diff)
		fmt.Fprintln(out, "!")
	}
}
