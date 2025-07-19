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
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		found := false
		for x := 0; x < 334 && !found; x++ {
			left := n - 3*x
			if left < 0 {
				break
			}
			z := 0
			for left > 0 && left%5 != 0 {
				z++
				left -= 7
			}
			if left >= 0 && left%5 == 0 {
				y := left / 5
				if 3*x+5*y+7*z == n {
					fmt.Fprintln(writer, x, y, z)
					found = true
					break
				}
			}
		}
		if !found {
			fmt.Fprintln(writer, -1)
		}
	}
}
