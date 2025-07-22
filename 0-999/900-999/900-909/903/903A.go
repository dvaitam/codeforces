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
	for ; n > 0; n-- {
		var x int
		fmt.Fscan(reader, &x)
		res := "NO"
		for a := 0; a <= x; a += 3 {
			if (x-a)%7 == 0 {
				res = "YES"
				break
			}
		}
		fmt.Fprintln(writer, res)
	}
}
