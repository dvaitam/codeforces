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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var L, R string
		fmt.Fscan(reader, &L, &R)
		// Pad with leading zeros to equal length
		if len(L) < len(R) {
			for len(L) < len(R) {
				L = "0" + L
			}
		} else if len(R) < len(L) {
			for len(R) < len(L) {
				R = "0" + R
			}
		}
		n := len(L)
		ans := 0
		for i := 0; i < n; i++ {
			if L[i] != R[i] {
				diff := int(L[i]) - int(R[i])
				if diff < 0 {
					diff = -diff
				}
				ans = diff + 9*(n-i-1)
				break
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
