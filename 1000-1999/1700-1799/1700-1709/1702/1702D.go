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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var w string
		fmt.Fscan(reader, &w)
		var p int
		fmt.Fscan(reader, &p)

		freq := make([]int, 26)
		total := 0
		for _, ch := range w {
			idx := int(ch - 'a')
			freq[idx]++
			total += idx + 1
		}

		keep := make([]int, 26)
		copy(keep, freq)

		for c := 25; c >= 0 && total > p; c-- {
			for keep[c] > 0 && total > p {
				keep[c]--
				total -= c + 1
			}
		}

		res := make([]byte, 0, len(w))
		for _, ch := range w {
			idx := int(ch - 'a')
			if keep[idx] > 0 {
				res = append(res, byte(ch))
				keep[idx]--
			}
		}
		fmt.Fprintln(writer, string(res))
	}
}
