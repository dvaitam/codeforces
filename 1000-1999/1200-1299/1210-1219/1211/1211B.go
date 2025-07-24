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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	var ans int64
	for i := 0; i < n; i++ {
		if a[i] == 0 {
			continue
		}
		candidate := int64(i+1) + int64(n)*(a[i]-1)
		if candidate > ans {
			ans = candidate
		}
	}

	fmt.Fprintln(writer, ans)
}
