package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	values := make(map[int64]int64, n)
	for i := 0; i < n; i++ {
		var a, x int64
		fmt.Fscan(reader, &a, &x)
		if x > values[a] {
			values[a] = x
		}
	}
	var m int
	fmt.Fscan(reader, &m)
	for i := 0; i < m; i++ {
		var b, y int64
		fmt.Fscan(reader, &b, &y)
		if y > values[b] {
			values[b] = y
		}
	}
	var sum int64
	for _, v := range values {
		sum += v
	}
	fmt.Println(sum)
}
