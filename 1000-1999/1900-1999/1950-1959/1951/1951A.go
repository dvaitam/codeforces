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
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)

		ones := 0
		positions := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				ones++
				positions = append(positions, i)
			}
		}

		result := "NO"
		if n <= 2 {
			if ones == 0 {
				result = "YES"
			}
		} else {
			if ones%2 == 0 {
				if ones == 2 {
					if positions[1]-positions[0] > 1 {
						result = "YES"
					}
				} else {
					result = "YES"
				}
			}
		}

		fmt.Fprintln(writer, result)
	}
}
