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

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	if len(s) == 0 {
		return
	}

	minChar := s[0]
	fmt.Fprintln(writer, "Mike")
	for i := 1; i < len(s); i++ {
		if s[i] > minChar {
			fmt.Fprintln(writer, "Ann")
		} else {
			fmt.Fprintln(writer, "Mike")
		}
		if s[i] < minChar {
			minChar = s[i]
		}
	}
}
