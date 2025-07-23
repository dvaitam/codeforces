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

	target := byte('a')
	b := []byte(s)
	for i := 0; i < len(b) && target <= 'z'; i++ {
		if b[i] <= target {
			b[i] = target
			target++
		}
	}

	if target <= 'z' {
		fmt.Fprintln(writer, "-1")
	} else {
		fmt.Fprintln(writer, string(b))
	}
}
