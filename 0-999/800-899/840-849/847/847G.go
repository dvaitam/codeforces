package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	counts := make([]int, 7)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j, c := range s {
			if c == '1' {
				counts[j]++
			}
		}
	}
	maxRooms := 0
	for _, v := range counts {
		if v > maxRooms {
			maxRooms = v
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, maxRooms)
	out.Flush()
}
