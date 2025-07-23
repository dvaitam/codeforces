package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ids[i])
	}
	for i := 1; i <= n; i++ {
		if k <= i {
			fmt.Println(ids[k-1])
			return
		}
		k -= i
	}
}
