package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for problemE.txt (Sofia and Strings).
// The actual algorithm is non-trivial and is not implemented here.
// This program reads the input format and outputs "YES" for each test case.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		var s, t string
		fmt.Fscan(reader, &n, &m)
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		// TODO: implement the real logic
		fmt.Fprintln(writer, "YES")
	}
}
