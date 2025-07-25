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

	// Compress consecutive identical characters into runs
	colors := make([]byte, 0, len(s))
	lens := make([]int, 0)
	for i := 0; i < len(s); {
		j := i
		for j < len(s) && s[j] == s[i] {
			j++
		}
		colors = append(colors, s[i])
		lens = append(lens, j-i)
		i = j
	}

	m := len(colors)
	if m%2 == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	mid := m / 2
	if lens[mid] < 2 {
		fmt.Fprintln(writer, 0)
		return
	}

	l, r := mid-1, mid+1
	for l >= 0 && r < m {
		if colors[l] != colors[r] || lens[l]+lens[r] < 3 {
			break
		}
		l--
		r++
	}

	if l < 0 && r >= m {
		fmt.Fprintln(writer, lens[mid]+1)
	} else {
		fmt.Fprintln(writer, 0)
	}
}
