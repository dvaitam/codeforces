package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0, 6)
	for i := 0; i < 6; i++ {
		if !in.Scan() {
			return
		}
		lines = append(lines, in.Text())
	}

	seatPositions := []int{0, 1, 3, 4, 6, 7}
	row, pos := -1, -1
	for i := 0; i < 6 && row == -1; i++ {
		for _, p := range seatPositions {
			if lines[i][p] == '.' {
				row, pos = i, p
				break
			}
		}
	}

	if row != -1 {
		b := []byte(lines[row])
		b[pos] = 'P'
		lines[row] = string(b)
	}

	out := bufio.NewWriter(os.Stdout)
	for _, l := range lines {
		fmt.Fprintln(out, l)
	}
	out.Flush()
}
