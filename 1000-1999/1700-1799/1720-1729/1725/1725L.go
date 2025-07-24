package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	// TODO: implement the actual algorithm. For now we only check a simple
	// necessary condition based on prefix and suffix sums. If either condition
	// fails we output -1, otherwise we output 0 as a placeholder.
	prefix := int64(0)
	for i := 0; i < n; i++ {
		prefix += a[i]
		if prefix < 0 {
			fmt.Println(-1)
			return
		}
	}
	suffix := int64(0)
	for i := n - 1; i >= 0; i-- {
		suffix += a[i]
		if suffix < 0 {
			fmt.Println(-1)
			return
		}
	}
	fmt.Println(0)
}
