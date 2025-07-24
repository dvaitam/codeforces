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
		var x int
		fmt.Fscan(reader, &x)
		doors := make([]int, 4)
		fmt.Fscan(reader, &doors[1], &doors[2], &doors[3])

		first := doors[x]
		if first == 0 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		second := doors[first]
		if second == 0 {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
