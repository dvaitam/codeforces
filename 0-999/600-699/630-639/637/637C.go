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
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &codes[i])
	}

	if n == 1 {
		fmt.Fprintln(writer, 6)
		return
	}

	minDist := 7
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := 0
			for k := 0; k < 6; k++ {
				if codes[i][k] != codes[j][k] {
					d++
				}
			}
			if d < minDist {
				minDist = d
				if minDist == 1 {
					break
				}
			}
		}
		if minDist == 1 {
			break
		}
	}

	if minDist < 1 {
		minDist = 1
	}
	k := (minDist - 1) / 2
	fmt.Fprintln(writer, k)
}
