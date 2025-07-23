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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Ints(arr)
	if arr[n-1] < arr[n] {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
