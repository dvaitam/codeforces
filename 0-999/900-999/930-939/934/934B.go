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

	var a int
	fmt.Fscan(reader, &a)
	if a > 36 {
		fmt.Fprintln(writer, -1)
		return
	}
	half := a / 2
	if a%2 == 0 {
		for i := 0; i < half; i++ {
			writer.WriteByte('8')
		}
		writer.WriteByte('\n')
	} else {
		for i := 0; i < half; i++ {
			writer.WriteByte('8')
		}
		writer.WriteByte('4')
		writer.WriteByte('\n')
	}
}
