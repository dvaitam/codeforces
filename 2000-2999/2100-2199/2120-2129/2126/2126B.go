package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		ans := 0
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		i := 0
		for i < n {
			if arr[i] == 1 {
				i++
				continue
			}
			j := i
			for j < n && arr[j] == 0 {
				j++
			}
			length := j - i
			ans += (length + 1) / (k + 1)
			i = j + 1
		}
		fmt.Fprintln(out, ans)
	}
}
