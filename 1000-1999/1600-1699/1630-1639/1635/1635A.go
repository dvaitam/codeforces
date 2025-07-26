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
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		res := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			res |= x
		}
		fmt.Fprintln(writer, res)
	}
}
