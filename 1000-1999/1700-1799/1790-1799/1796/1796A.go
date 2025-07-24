package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	base := "FBFFBFFB"
	pattern := strings.Repeat(base, 20)

	for ; t > 0; t-- {
		var k int
		var s string
		fmt.Fscan(reader, &k, &s)
		if strings.Contains(pattern, s) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
