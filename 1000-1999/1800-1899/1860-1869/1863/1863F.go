package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement a correct solution for problem F.
// The current implementation is a placeholder that always
// outputs a string of '0's for each test case.

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		for i := 0; i < n; i++ {
			out.WriteByte('0')
		}
		out.WriteByte('\n')
	}
}
