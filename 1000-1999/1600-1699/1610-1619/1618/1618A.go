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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		arr := make([]int, 7)
		for i := 0; i < 7; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		a1 := arr[0]
		a2 := arr[1]
		a3 := arr[6] - a1 - a2
		fmt.Fprintln(writer, a1, a2, a3)
	}
}
