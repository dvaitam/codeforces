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
	var minV int64 = 1<<63 - 1
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] < minV {
			minV = a[i]
		}
	}

	longest := 0
	cur := 0
	for i := 0; i < 2*n; i++ {
		if a[i%n] > minV {
			cur++
			if cur > longest {
				longest = cur
			}
		} else {
			if cur > longest {
				longest = cur
			}
			cur = 0
		}
	}
	if cur > longest {
		longest = cur
	}

	result := int64(n)*minV + int64(longest)
	fmt.Fprintln(writer, result)
}
