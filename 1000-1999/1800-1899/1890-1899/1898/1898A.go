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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		cnt := 0
		for i := 0; i < n; i++ {
			if s[i] == 'B' {
				cnt++
			}
		}
		if cnt == k {
			fmt.Fprintln(out, 0)
			continue
		}
		if cnt < k {
			diff := k - cnt
			prefA := 0
			idx := 0
			for i := 0; i < n; i++ {
				if s[i] == 'A' {
					prefA++
				}
				if prefA == diff {
					idx = i + 1
					break
				}
			}
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, idx, "B")
		} else {
			diff := cnt - k
			prefB := 0
			idx := 0
			for i := 0; i < n; i++ {
				if s[i] == 'B' {
					prefB++
				}
				if prefB == diff {
					idx = i + 1
					break
				}
			}
			fmt.Fprintln(out, 1)
			fmt.Fprintln(out, idx, "A")
		}
	}
}
