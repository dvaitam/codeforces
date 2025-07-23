package main

import (
	"bufio"
	"fmt"
	"os"
)

func normalize(x int) int {
	for x%2 == 0 {
		x /= 2
	}
	for x%3 == 0 {
		x /= 3
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	base := normalize(arr[0])
	for i := 1; i < n; i++ {
		if normalize(arr[i]) != base {
			fmt.Fprintln(out, "No")
			return
		}
	}
	fmt.Fprintln(out, "Yes")
}
