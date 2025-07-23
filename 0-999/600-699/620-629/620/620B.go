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

	var a, b int
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}

	seg := []int{6, 2, 5, 5, 4, 5, 6, 3, 7, 6}
	total := 0
	for i := a; i <= b; i++ {
		x := i
		if x == 0 {
			total += seg[0]
		}
		for x > 0 {
			total += seg[x%10]
			x /= 10
		}
	}
	fmt.Fprintln(writer, total)
}
