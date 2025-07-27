package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		evens := make([]int, 0)
		odds := make([]int, 0)
		for _, v := range arr {
			if v%2 == 0 {
				evens = append(evens, v)
			} else {
				odds = append(odds, v)
			}
		}
		b := append(evens, odds...)
		count := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if gcd(b[i], 2*b[j]) > 1 {
					count++
				}
			}
		}
		fmt.Fprintln(writer, count)
	}
}
