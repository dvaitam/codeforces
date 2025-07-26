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
		var n int
		fmt.Fscan(in, &n)
		if n == 1 {
			fmt.Fprintln(out, "a")
			continue
		}
		k := 2
		if n != 2 {
			for n%k == 0 {
				k++
			}
		}
		res := make([]byte, n)
		for i := 0; i < n; i++ {
			res[i] = byte('a' + i%k)
		}
		fmt.Fprintln(out, string(res))
	}
}
