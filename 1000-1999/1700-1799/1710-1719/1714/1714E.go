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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if solve(arr) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}

func solve(arr []int) bool {
	hasZero := false
	for i, v := range arr {
		if v%10 == 5 {
			v += 5
		}
		if v%10 == 0 {
			hasZero = true
		}
		arr[i] = v
	}
	if hasZero {
		first := arr[0]
		for _, v := range arr {
			if v != first {
				return false
			}
		}
		return true
	}
	parity := -1
	for i, v := range arr {
		for v%10 != 2 {
			v += v % 10
		}
		arr[i] = v
		p := (v / 10) % 2
		if parity == -1 {
			parity = p
		} else if parity != p {
			return false
		}
	}
	return true
}
