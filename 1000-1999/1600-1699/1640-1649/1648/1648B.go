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
		var n, c int
		fmt.Fscan(reader, &n, &c)
		freq := make([]int, c+1)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if x <= c {
				freq[x]++
			}
		}
		prefix := make([]int, c+1)
		for i := 1; i <= c; i++ {
			prefix[i] = prefix[i-1] + freq[i]
		}
		ok := true
		for y := 1; y <= c && ok; y++ {
			if freq[y] == 0 {
				continue
			}
			for k := 1; k*y <= c && ok; k++ {
				L := k * y
				R := (k+1)*y - 1
				if R > c {
					R = c
				}
				if prefix[R]-prefix[L-1] > 0 {
					if freq[k] == 0 {
						ok = false
						break
					}
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
