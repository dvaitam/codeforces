package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the game described in problemA.txt for contest 1673.
// It determines the winner and the difference between the players' scores.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		total := 0
		for i := 0; i < n; i++ {
			total += int(s[i]-'a') + 1
		}
		if n == 1 {
			fmt.Fprintln(writer, "Bob", total)
			continue
		}
		if n%2 == 0 {
			fmt.Fprintln(writer, "Alice", total)
			continue
		}
		// n is odd and > 1
		first := int(s[0]-'a') + 1
		last := int(s[n-1]-'a') + 1
		bob := first
		if last < bob {
			bob = last
		}
		alice := total - bob
		if alice > bob {
			fmt.Fprintln(writer, "Alice", alice-bob)
		} else {
			fmt.Fprintln(writer, "Bob", bob-alice)
		}
	}
}
