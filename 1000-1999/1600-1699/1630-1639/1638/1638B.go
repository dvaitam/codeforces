package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		lastOdd := -1
		lastEven := -1
		ok := true
		for _, v := range arr {
			if v%2 == 0 {
				if lastEven != -1 && v < lastEven {
					ok = false
					break
				}
				lastEven = v
			} else {
				if lastOdd != -1 && v < lastOdd {
					ok = false
					break
				}
				lastOdd = v
			}
		}
		if ok {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
