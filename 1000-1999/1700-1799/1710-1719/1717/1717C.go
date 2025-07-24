package main

import (
	"bufio"
	"fmt"
	"os"
)

func canTransform(a, b []int64) bool {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i] > b[i] {
			return false
		}
	}
	same := true
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			same = false
			break
		}
	}
	if same {
		return true
	}
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if b[i] > b[j]+1 && a[i] != b[i] {
			return false
		}
	}
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if b[i] <= b[j]+1 {
			return true
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		if canTransform(a, b) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
