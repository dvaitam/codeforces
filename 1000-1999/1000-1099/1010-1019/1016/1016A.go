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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	blank := m
	for i := 0; i < n; i++ {
		var ai int
		fmt.Fscan(reader, &ai)
		turn := 0
		if ai <= blank {
			blank -= ai
		} else {
			ai -= blank
			used := (ai + m - 1) / m
			blank = used*m - ai
			turn = used
		}
		if blank == 0 {
			turn++
			blank = m
		}
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, turn)
	}
	fmt.Fprint(writer, "\n")
}
