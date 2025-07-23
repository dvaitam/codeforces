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

	values := make([]int, 0)
	pos := 0
	for n > 0 {
		if n&1 == 1 {
			values = append(values, pos+1)
		}
		n >>= 1
		pos++
	}

	for i := len(values) - 1; i >= 0; i-- {
		if i != len(values)-1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, values[i])
	}
	fmt.Fprintln(writer)
}
