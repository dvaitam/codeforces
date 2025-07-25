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
		var n int
		fmt.Fscan(in, &n)
		found := false
		for i := 1; i <= 26 && !found; i++ {
			for j := 1; j <= 26 && !found; j++ {
				k := n - i - j
				if k >= 1 && k <= 26 {
					fmt.Fprintf(out, "%c%c%c\n", 'a'+i-1, 'a'+j-1, 'a'+k-1)
					found = true
				}
			}
		}
	}
}
