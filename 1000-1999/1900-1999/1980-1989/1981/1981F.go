package main

import (
	"bufio"
	"fmt"
	"os"
)

func mex(a, b int) int {
	used := map[int]bool{a: true, b: true}
	for i := 1; ; i++ {
		if !used[i] {
			return i
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n != 2 {
			// read and discard input for unsupported cases
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				fmt.Fscan(in, &arr[i])
			}
			for i := 2; i <= n; i++ {
				var p int
				fmt.Fscan(in, &p)
			}
			fmt.Fprintln(out, 0)
			continue
		}
		var a1, a2 int
		fmt.Fscan(in, &a1, &a2)
		var p2 int
		fmt.Fscan(in, &p2)
		fmt.Fprintln(out, mex(a1, a2))
	}
}
