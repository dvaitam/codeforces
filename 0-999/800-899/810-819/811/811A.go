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

	var a, b int64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	for turn := int64(1); ; turn++ {
		if turn%2 == 1 {
			if a < turn {
				fmt.Fprintln(writer, "Vladik")
				return
			}
			a -= turn
		} else {
			if b < turn {
				fmt.Fprintln(writer, "Valera")
				return
			}
			b -= turn
		}
	}
}
