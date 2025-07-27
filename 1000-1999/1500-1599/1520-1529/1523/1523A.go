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
		var m int
		fmt.Fscan(reader, &n, &m)
		var s string
		fmt.Fscan(reader, &s)
		a := []byte(s)
		iterations := m
		if iterations > n {
			iterations = n
		}
		for iter := 0; iter < iterations; iter++ {
			changed := false
			b := make([]byte, n)
			copy(b, a)
			for i := 0; i < n; i++ {
				if a[i] == '0' {
					cnt := 0
					if i > 0 && a[i-1] == '1' {
						cnt++
					}
					if i+1 < n && a[i+1] == '1' {
						cnt++
					}
					if cnt == 1 {
						b[i] = '1'
						changed = true
					} else {
						b[i] = '0'
					}
				} else {
					b[i] = '1'
				}
			}
			if !changed {
				break
			}
			a = b
		}
		fmt.Fprintln(writer, string(a))
	}
}
