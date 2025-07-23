package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(r, g, b int) string {
	if r > 0 && g > 0 && b > 0 {
		return "BGR"
	}
	if g == 0 && b == 0 {
		return "R"
	}
	if r == 0 && b == 0 {
		return "G"
	}
	if r == 0 && g == 0 {
		return "B"
	}
	if r == 0 {
		// colors G and B present
		if g == 1 && b == 1 {
			return "R"
		}
		if g == 1 && b > 1 {
			return "GR"
		}
		if g > 1 && b == 1 {
			return "BR"
		}
		return "BGR"
	}
	if g == 0 {
		// colors R and B present
		if r == 1 && b == 1 {
			return "G"
		}
		if r == 1 && b > 1 {
			return "GR"
		}
		if r > 1 && b == 1 {
			return "BG"
		}
		return "BGR"
	}
	// b == 0, colors R and G present
	if r == 1 && g == 1 {
		return "B"
	}
	if r == 1 && g > 1 {
		return "BR"
	}
	if r > 1 && g == 1 {
		return "BG"
	}
	return "BGR"
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	fmt.Fscan(reader, &s)

	r, g, b := 0, 0, 0
	for i := 0; i < n; i++ {
		switch s[i] {
		case 'R':
			r++
		case 'G':
			g++
		case 'B':
			b++
		}
	}

	result := solve(r, g, b)
	fmt.Fprintln(writer, result)
}
