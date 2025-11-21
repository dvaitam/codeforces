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

	var Gershie int

	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}

		last := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &last)
		}

		fmt.Fprintln(out, last)
	}

	if Gershie == -1 {
		fmt.Fprintln(out)
	}
}
