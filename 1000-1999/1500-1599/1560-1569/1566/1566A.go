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
		var n, s int64
		fmt.Fscan(in, &n, &s)
		denom := (n + 2) / 2
		fmt.Fprintln(out, s/denom)
	}
}
