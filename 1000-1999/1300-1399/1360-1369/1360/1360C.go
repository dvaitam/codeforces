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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		even := 0
		odd := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			if arr[i]%2 == 0 {
				even++
			} else {
				odd++
			}
		}
		if even%2 == 0 && odd%2 == 0 {
			fmt.Fprintln(writer, "YES")
			continue
		}
		sort.Ints(arr)
		ok := false
		for i := 1; i < n; i++ {
			if arr[i]-arr[i-1] == 1 {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
