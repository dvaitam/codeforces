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
	if a < 0 {
		return -a
	}
	return a
}

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

	g := 0
	for _, v := range arr {
		g = gcd(g, v)
	}
	if g > 1 {
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, 0)
		return
	}

	ops := 0
	for i := 0; i < n; i++ {
		if arr[i]%2 != 0 {
			if i+1 < n {
				if arr[i+1]%2 != 0 {
					ops++
					i++
				} else {
					ops += 2
					i++
				}
			} else {
				ops += 2
			}
		}
	}

	fmt.Fprintln(writer, "YES")
	fmt.Fprintln(writer, ops)
}
