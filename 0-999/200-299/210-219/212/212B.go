package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	var m int
	fmt.Fscan(in, &m)

	for ; m > 0; m-- {
		var t string
		fmt.Fscan(in, &t)
		mask := 0
		for _, ch := range t {
			mask |= 1 << (ch - 'a')
		}
		count := 0
		i := 0
		for i < n {
			if mask&(1<<(s[i]-'a')) == 0 {
				i++
				continue
			}
			seen := 0
			j := i
			for j < n {
				bit := 1 << (s[j] - 'a')
				if mask&bit == 0 {
					break
				}
				seen |= bit
				j++
			}
			if seen == mask {
				count++
			}
			i = j
		}
		fmt.Fprintln(out, count)
	}
}
