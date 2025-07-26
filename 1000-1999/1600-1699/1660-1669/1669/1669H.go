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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		result := 0
		for b := 30; b >= 0; b-- {
			count := 0
			for i := 0; i < n; i++ {
				if (arr[i]>>uint(b))&1 == 1 {
					count++
				}
			}
			need := n - count
			if need <= k {
				k -= need
				result |= 1 << uint(b)
			}
		}
		fmt.Fprintln(writer, result)
	}
}
