package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSorted(a []int) bool {
	if len(a) == 0 {
		return true
	}
	if a[0] == 0 {
		return false
	}
	for i := 1; i < len(a); i++ {
		if a[i] == 0 || a[i-1] > a[i] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, c int
	if _, err := fmt.Fscan(reader, &n, &m, &c); err != nil {
		return
	}

	arr := make([]int, n)
	half := (c + 1) / 2
	for round := 0; round < m; round++ {
		var x int
		if _, err := fmt.Fscan(reader, &x); err != nil {
			return
		}
		var pos int
		if x <= half {
			pos = -1
			for i := 0; i < n; i++ {
				if arr[i] == 0 || arr[i] > x {
					pos = i
					break
				}
			}
		} else {
			pos = -1
			for i := n - 1; i >= 0; i-- {
				if arr[i] == 0 || arr[i] < x {
					pos = i
					break
				}
			}
		}
		if pos == -1 {
			pos = 0
		}
		arr[pos] = x
		fmt.Fprintln(writer, pos+1)
		writer.Flush()
		if isSorted(arr) {
			return
		}
	}
}
