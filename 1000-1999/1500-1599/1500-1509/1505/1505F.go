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

	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}

	if a > 0 {
		fmt.Fprintln(out, 1)
	} else if a == 0 {
		fmt.Fprintln(out, 0)
	} else {
		fmt.Fprintln(out, -1)
	}
}
