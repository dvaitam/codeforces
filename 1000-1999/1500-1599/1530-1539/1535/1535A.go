package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s1, s2, s3, s4 int
		if _, err := fmt.Fscan(in, &s1, &s2, &s3, &s4); err != nil {
			return
		}
		arr := []int{s1, s2, s3, s4}
		sort.Ints(arr)
		top1 := max(s1, s2)
		top2 := max(s3, s4)
		if (top1 == arr[3] && top2 == arr[2]) || (top1 == arr[2] && top2 == arr[3]) {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
