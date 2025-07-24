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
		n := len(s)
		zero := make([]int, n)
		for i := 1; i < n; i++ {
			if s[i-1] == '0' {
				zero[i] = zero[i-1] + 1
			} else {
				zero[i] = 0
			}
		}
		ans := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				val := 0
				for j := i; j < n && j-i < 20; j++ {
					val = val*2 + int(s[j]-'0')
					length := j - i + 1
					if val >= length {
						need := val - length
						if need <= zero[i] {
							ans++
						}
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
