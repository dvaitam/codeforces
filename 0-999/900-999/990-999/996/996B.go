package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	var ans int
	minTime := int64(1<<63 - 1)
	nn := int64(n)
	for i := 0; i < n; i++ {
		t := int64(i)
		if a[i] > int64(i) {
			diff := a[i] - int64(i)
			cycles := (diff + nn - 1) / nn
			t += cycles * nn
		}
		if t < minTime {
			minTime = t
			ans = i
		}
	}
	fmt.Println(ans + 1)
}
