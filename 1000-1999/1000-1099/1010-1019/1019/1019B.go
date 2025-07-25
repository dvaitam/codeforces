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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	res := -1
	half := n / 2
	for i := 0; i < half; i++ {
		if a[i] == a[i+half] {
			res = i + 1
			break
		}
	}
	fmt.Fprintln(writer, res)
}
