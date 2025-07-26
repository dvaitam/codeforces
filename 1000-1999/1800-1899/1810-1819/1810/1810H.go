package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement the algorithm described in problemH.txt.
// This placeholder parses input but outputs 0 for each test case.
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
		fmt.Fprintln(out, 0)
	}
}
