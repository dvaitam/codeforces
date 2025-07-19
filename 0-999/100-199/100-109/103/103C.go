package main

import (
	"bufio"
	"fmt"
	"os"
)

func cal(n, k, x int64) byte {
	// Handle odd initial n case
	if n&1 == 1 {
		if x == n {
			if k > 0 {
				return 'X'
			}
			return '.'
		}
		n--
		k--
	}
	if x%2 == 0 {
		if x <= n-k*2 {
			return '.'
		}
		return 'X'
	}
	if x < n-(k-n/2)*2 {
		return '.'
	}
	return 'X'
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int64
	var p int
	if _, err := fmt.Fscan(reader, &n, &k, &p); err != nil {
		return
	}
	result := make([]byte, p)
	for i := 0; i < p; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		result[i] = cal(n, k, x)
	}
	writer.Write(result)
	writer.WriteByte('\n')
}
