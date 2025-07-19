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
	fmt.Fscan(reader, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var c int
	for i := 0; i < n; i++ {
		if arr[i] == 0 {
			if i == 0 || i == n-1 {
				continue
			} else if arr[i-1] == 1 && arr[i+1] == 1 {
				arr[i+1] = 0
				c++
			}
		}
	}
	fmt.Fprint(writer, c)
}
