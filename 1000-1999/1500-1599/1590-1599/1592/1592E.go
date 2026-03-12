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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	ans := 0
	for l := 0; l < n; l++ {
		andVal := arr[l]
		xorVal := arr[l]
		for r := l + 1; r < n; r++ {
			andVal &= arr[r]
			xorVal ^= arr[r]
			if andVal > xorVal {
				length := r - l + 1
				if length > ans {
					ans = length
				}
			}
			if andVal == 0 {
				break
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
