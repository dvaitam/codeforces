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

	var sum int64
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(in, &x)
		sum += x
	}

	fmt.Fprintln(out, sum)
}
