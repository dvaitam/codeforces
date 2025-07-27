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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			freq[arr[i]]++
		}
		uniqueVal := 0
		for v, c := range freq {
			if c == 1 {
				uniqueVal = v
				break
			}
		}
		for i, v := range arr {
			if v == uniqueVal {
				fmt.Fprintln(writer, i+1)
				break
			}
		}
	}
}
