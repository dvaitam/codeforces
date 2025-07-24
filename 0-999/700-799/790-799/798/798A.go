package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	diff := 0
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			diff++
		}
	}
	if diff == 1 || (diff == 0 && n%2 == 1) {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
