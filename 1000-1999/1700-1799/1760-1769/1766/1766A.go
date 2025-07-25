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
		pow := 1
		for pow*10 <= n {
			pow *= 10
		}
		digits := 0
		for tmp := n; tmp > 0; tmp /= 10 {
			digits++
		}
		ans := 9*(digits-1) + n/pow
		fmt.Fprintln(writer, ans)
	}
}
