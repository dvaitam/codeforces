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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, d int
		fmt.Fscan(reader, &n, &d)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if canConstruct(arr, d) {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}

func canConstruct(a []int, d int) bool {
	mandatory := 0
	rest := make([]int, 0, len(a))
	for _, x := range a {
		if x*2 > d {
			mandatory += x
		} else {
			rest = append(rest, x)
		}
	}
	if mandatory > d {
		return false
	}
	target := d - mandatory
	dp := make([]bool, target+1)
	dp[0] = true
	for _, x := range rest {
		for s := target; s >= x; s-- {
			if dp[s-x] {
				dp[s] = true
			}
		}
	}
	return dp[target]
}
