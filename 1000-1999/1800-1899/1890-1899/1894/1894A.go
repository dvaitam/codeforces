package main

import (
	"bufio"
	"fmt"
	"os"
)

func winnerFor(X, Y int, s string) (byte, bool) {
	setsA, setsB := 0, 0
	winsA, winsB := 0, 0
	n := len(s)
	for i := 0; i < n; i++ {
		if setsA == Y || setsB == Y {
			// game would have ended earlier
			return 0, false
		}
		if s[i] == 'A' {
			winsA++
		} else {
			winsB++
		}
		if winsA == X {
			setsA++
			winsA, winsB = 0, 0
			if setsA == Y {
				if i != n-1 {
					return 0, false
				}
				return 'A', true
			}
		} else if winsB == X {
			setsB++
			winsA, winsB = 0, 0
			if setsB == Y {
				if i != n-1 {
					return 0, false
				}
				return 'B', true
			}
		}
	}
	return 0, false
}

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
		winners := make(map[byte]struct{})
		for x := 1; x <= n; x++ {
			for y := 1; y <= n; y++ {
				w, ok := winnerFor(x, y, s)
				if ok {
					winners[w] = struct{}{}
				}
			}
		}
		if len(winners) == 1 {
			for k := range winners {
				fmt.Fprintf(writer, "%c\n", k)
			}
		} else {
			fmt.Fprintln(writer, "?")
		}
	}
}
