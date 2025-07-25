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
	fmt.Fscan(reader, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		val := a[i] + int64(n-i-1)
		if prefix[i] > val {
			prefix[i+1] = prefix[i]
		} else {
			prefix[i+1] = val
		}
	}

	suffix := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		val := a[i] + int64(i)
		if suffix[i+1] > val {
			suffix[i] = suffix[i+1]
		} else {
			suffix[i] = val
		}
	}

	ans := int64(1<<63 - 1)
	for i := 0; i < n; i++ {
		val := a[i]
		if prefix[i] > val {
			val = prefix[i]
		}
		if suffix[i+1] > val {
			val = suffix[i+1]
		}
		if val < ans {
			ans = val
		}
	}

	fmt.Fprintln(writer, ans)
}
