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

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	for a := int64(0); a*1234567 <= n; a++ {
		remA := n - a*1234567
		for b := int64(0); b*123456 <= remA; b++ {
			rem := remA - b*123456
			if rem%1234 == 0 {
				fmt.Fprintln(out, "YES")
				return
			}
		}
	}
	fmt.Fprintln(out, "NO")
}
