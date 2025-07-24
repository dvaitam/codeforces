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

	var m, d int
	if _, err := fmt.Fscan(reader, &m, &d); err != nil {
		return
	}

	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	total := days[m-1] + d - 1
	columns := (total + 6) / 7

	fmt.Fprintln(writer, columns)
}
