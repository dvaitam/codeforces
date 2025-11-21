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
	var mode string
	fmt.Fscan(in, &t)
	fmt.Fscan(in, &mode) // expect "manual"

	for i := 0; i < t; i++ {
		var W int
		fmt.Fscan(in, &W)
		fmt.Fprintln(out, W)
	}
}
