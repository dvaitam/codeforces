package main

import (
	"bufio"
	"fmt"
	"os"
)

const piDigits = "314159265358979323846264338327"

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		count := 0
		for i := 0; i < len(s) && i < len(piDigits); i++ {
			if s[i] == piDigits[i] {
				count++
			} else {
				break
			}
		}
		fmt.Fprintln(writer, count)
	}
}
