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
	for tc := 0; tc < t; tc++ {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		c := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				c++
			}
		}
		if c == 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		// find smallest power of two >= c
		target := 1
		for target < c {
			target <<= 1
		}
		cuts := make([]bool, n)
		cuts[n-1] = true
		deficit := target - c
		v := 0
		for i := 0; i < n-1; i++ {
			if s[i] == '1' {
				v++
			}
			if v <= deficit {
				deficit -= v
				v <<= 1
			} else {
				// make a cut here
				cuts[i] = true
				v = 0
			}
		}
		// output segments
		// deficit should be zero
		count := 0
		for i := 0; i < n; i++ {
			if cuts[i] {
				count++
			}
		}
		fmt.Fprintln(writer, count)
		l := 1
		for i := 0; i < n; i++ {
			if cuts[i] {
				fmt.Fprintln(writer, l, i+1)
				l = i + 2
			}
		}
	}
}
