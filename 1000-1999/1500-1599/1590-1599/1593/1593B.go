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
		ans := len(s) // upper bound
		targets := []string{"00", "25", "50", "75"}
		for _, p := range targets {
			p0 := p[0]
			p1 := p[1]
			pos1 := -1
			// find last char p1 from the end
			for i := len(s) - 1; i >= 0; i-- {
				if s[i] == p1 {
					pos1 = i
					break
				}
			}
			if pos1 == -1 {
				continue
			}
			pos0 := -1
			for i := pos1 - 1; i >= 0; i-- {
				if s[i] == p0 {
					pos0 = i
					break
				}
			}
			if pos0 == -1 {
				continue
			}
			moves := (len(s) - pos1 - 1) + (pos1 - pos0 - 1)
			if moves < ans {
				ans = moves
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
