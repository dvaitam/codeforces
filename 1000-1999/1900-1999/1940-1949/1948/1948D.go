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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		b := []byte(s)
		n := len(b)
		result := 0
		for k := n / 2; k >= 1; k-- {
			run := 0
			found := false
			for i := n - k - 1; i >= 0; i-- {
				c1 := b[i]
				c2 := b[i+k]
				if c1 == c2 || c1 == '?' || c2 == '?' {
					run++
					if run >= k {
						result = 2 * k
						found = true
						break
					}
				} else {
					run = 0
				}
			}
			if found {
				break
			}
		}
		fmt.Fprintln(writer, result)
	}
}
