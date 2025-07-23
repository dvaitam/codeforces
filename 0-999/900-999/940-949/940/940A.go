package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, d int
	if _, err := fmt.Fscan(reader, &n, &d); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	sort.Ints(arr)
	left := 0
	best := 0
	for right := 0; right < n; right++ {
		for arr[right]-arr[left] > d {
			left++
		}
		if right-left+1 > best {
			best = right - left + 1
		}
	}
	fmt.Fprintln(writer, n-best)
}
