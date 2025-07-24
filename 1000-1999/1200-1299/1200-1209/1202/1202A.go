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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var x, y string
		fmt.Fscan(in, &x)
		fmt.Fscan(in, &y)

		// find index of the least significant '1' in y
		trailing := 0
		for i := len(y) - 1; i >= 0 && y[i] == '0'; i-- {
			trailing++
		}

		// find first '1' in x at or after position 'trailing' from the right
		j := trailing
		for j < len(x) && x[len(x)-1-j] == '0' {
			j++
		}
		fmt.Fprintln(out, j-trailing)
	}
}
