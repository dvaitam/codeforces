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
		var n int
		var k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		ones := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				ones++
			}
		}
		if ones == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		arr := []byte(s)
		first := int(arr[0] - '0')
		last := int(arr[n-1] - '0')

		if last == 0 {
			for i := n - 2; i >= 0; i-- {
				if arr[i] == '1' {
					dist := (n - 1) - i
					if dist <= k {
						k -= dist
						arr[i] = '0'
						arr[n-1] = '1'
						last = 1
						if i == 0 {
							first = 0
						}
					}
					break
				}
			}
		}
		if first == 0 {
			for i := 0; i < n; i++ {
				if arr[i] == '1' {
					dist := i
					if dist <= k {
						k -= dist
						arr[i] = '0'
						arr[0] = '1'
						first = 1
						if i == n-1 {
							last = 0
						}
					}
					break
				}
			}
		}

		result := 11*ones - 10*last - first
		fmt.Fprintln(writer, result)
	}
}
