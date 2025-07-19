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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for t > 0 {
		t--
		var n int
		fmt.Fscan(in, &n)
		if n%7 == 0 {
			fmt.Fprintln(out, n)
		} else {
			base := n - n%10
			for i := 0; i < 10; i++ {
				if (base+i)%7 == 0 {
					fmt.Fprintln(out, base+i)
					break
				}
			}
		}
	}
}
