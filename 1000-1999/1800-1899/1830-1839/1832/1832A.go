package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		freq := make(map[rune]int)
		for _, c := range s {
			freq[c]++
		}
		if len(freq) == 1 {
			fmt.Fprintln(out, "NO")
			continue
		}
		if len(freq) == 2 {
			ones := 0
			for _, c := range freq {
				if c == 1 {
					ones++
				}
			}
			if ones == 1 {
				fmt.Fprintln(out, "NO")
				continue
			}
		}
		fmt.Fprintln(out, "YES")
	}
}
