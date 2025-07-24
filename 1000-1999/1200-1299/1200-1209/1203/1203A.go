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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		idx := -1
		for i, v := range arr {
			if v == 1 {
				idx = i
				break
			}
		}
		ok1 := true
		for i := 0; i < n; i++ {
			if arr[(idx+i)%n] != i+1 {
				ok1 = false
				break
			}
		}
		ok2 := true
		for i := 0; i < n; i++ {
			// counterclockwise: decreasing order
			if arr[(idx-i+n)%n] != i+1 {
				ok2 = false
				break
			}
		}
		if ok1 || ok2 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
