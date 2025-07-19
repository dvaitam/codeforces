package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	data, _ := io.ReadAll(reader)
	idx := 0
	n, idx := readInt(data, idx)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], idx = readInt(data, idx)
	}
	x := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		if a[i] != 0 {
			x[i] = x[i+1] + 1
		} else {
			x[i] = x[i+1]
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		ones := x[i]
		// suffix length = n - i
		if ones == 0 || ones == n-i {
			fmt.Fprintln(writer, i)
			return
		}
	}
	// fallback
	fmt.Fprintln(writer, n)
}

// readInt parses an integer from data starting at idx, returning the integer and new index
func readInt(b []byte, idx int) (int, int) {
	n := len(b)
	// skip non-numeric
	for idx < n && (b[idx] < '0' || b[idx] > '9') && b[idx] != '-' {
		idx++
	}
	neg := false
	if idx < n && b[idx] == '-' {
		neg = true
		idx++
	}
	val := 0
	for idx < n && b[idx] >= '0' && b[idx] <= '9' {
		val = val*10 + int(b[idx]-'0')
		idx++
	}
	if neg {
		val = -val
	}
	return val, idx
}
