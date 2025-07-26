package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves problemA.txt for contest 1684A (Digit Minimization).
// Alice can rearrange the digits before each removal, effectively choosing
// which digit survives. Therefore the smallest possible final digit is the
// minimum digit present in the given number.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n string
		fmt.Fscan(reader, &n)
		minDigit := byte('9')
		for i := 0; i < len(n); i++ {
			if n[i] < minDigit {
				minDigit = n[i]
			}
		}
		fmt.Fprintln(writer, string(minDigit))
	}
}
