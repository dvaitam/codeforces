package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		var token string
		fmt.Fscan(in, &token)

		arr := make([]int, n)
		idx := 0

		if token == "manual" {
			for i := 0; i < n; i++ {
				fmt.Fscan(in, &arr[i])
			}
		} else {
			val, err := strconv.Atoi(token)
			if err != nil {
				val = 0
			}
			arr[idx] = val
			idx++
			for i := idx; i < n; i++ {
				fmt.Fscan(in, &arr[i])
			}
		}

		pos := -1
		for i, v := range arr {
			if v == -1 {
				pos = i + 1
				break
			}
		}
		if pos == -1 {
			pos = 1
		}
		fmt.Fprintln(out, pos)
	}
}
