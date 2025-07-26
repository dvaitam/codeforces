package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt for contest 1611A.
// We can only reverse prefixes of the number. To make the number even, we just
// need any even digit to become the last digit. If the last digit is already
// even, no operation is needed. If the first digit is even, reversing the whole
// number puts it at the end in one move. Otherwise, if any digit inside is even,
// we can move it to the front then reverse the whole number, requiring two
// moves. If there is no even digit at all, it's impossible.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		n := len(s)
		lastEven := (s[n-1]-'0')%2 == 0
		firstEven := (s[0]-'0')%2 == 0
		if lastEven {
			fmt.Fprintln(out, 0)
			continue
		}
		if firstEven {
			fmt.Fprintln(out, 1)
			continue
		}
		found := false
		for i := 1; i < n-1; i++ {
			if (s[i]-'0')%2 == 0 {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(out, 2)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
