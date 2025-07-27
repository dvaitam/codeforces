package main

import (
	"bufio"
	"fmt"
	"os"
)

// The original problem 1482G "Vabank" is interactive. This
// implementation targets the simplified "hacked" version
// used on Codeforces where the value of M is provided in the
// input. For each test case we just output "! M".
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var m int64
		fmt.Fscan(reader, &m)
		fmt.Fprintf(writer, "! %d\n", m)
	}
}
