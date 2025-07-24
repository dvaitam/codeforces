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
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		days := 1
		present := make(map[byte]bool)
		for i := 0; i < len(s); i++ {
			c := s[i]
			if !present[c] {
				if len(present) == 3 {
					days++
					for k := range present {
						delete(present, k)
					}
				}
				present[c] = true
			}
		}
		fmt.Fprintln(writer, days)
	}
}
