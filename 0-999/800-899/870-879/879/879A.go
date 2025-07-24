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
	last := 0
	for i := 0; i < n; i++ {
		var s, d int
		fmt.Fscan(reader, &s, &d)
		day := s
		for day <= last {
			day += d
		}
		last = day
	}
	fmt.Fprint(writer, last)
}
