package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var hashCount int64
		zeroExists := false
		for i := 0; i < n; i++ {
			var row string
			fmt.Fscan(reader, &row)
			for _, ch := range row {
				if ch == '#' {
					hashCount++
				} else if ch == '0' {
					zeroExists = true
				}
			}
		}
		ans := modPow(2, hashCount)
		if !zeroExists {
			ans--
			if ans < 0 {
				ans += mod
			}
		}
		fmt.Fprintln(writer, ans%mod)
	}
}
