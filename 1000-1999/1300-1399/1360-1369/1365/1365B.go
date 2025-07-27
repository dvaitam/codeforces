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
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		has0, has1 := false, false
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
			if b[i] == 0 {
				has0 = true
			} else {
				has1 = true
			}
		}
		sorted := true
		for i := 1; i < n; i++ {
			if a[i] < a[i-1] {
				sorted = false
				break
			}
		}
		if sorted || (has0 && has1) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
