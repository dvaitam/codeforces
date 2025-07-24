package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Ints(arr)
	used := make([]bool, n)
	var colors int
	for i := 0; i < n; i++ {
		if used[i] {
			continue
		}
		colors++
		for j := i; j < n; j++ {
			if !used[j] && arr[j]%arr[i] == 0 {
				used[j] = true
			}
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprint(writer, colors)
}
