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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		liked := make([]bool, n)
		likedCnt := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				liked[i] = true
				likedCnt++
			}
		}
		if likedCnt == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		idx := 0
		eaten := 0
		for likedCnt > 0 {
			cur := arr[idx]
			if liked[cur] {
				liked[cur] = false
				likedCnt--
			}
			arr = append(arr[:idx], arr[idx+1:]...)
			eaten++
			if likedCnt == 0 || len(arr) == 0 {
				break
			}
			if idx == len(arr) {
				idx = 0
			}
			idx = (idx + k - 1) % len(arr)
		}
		fmt.Fprintln(out, eaten)
	}
}
