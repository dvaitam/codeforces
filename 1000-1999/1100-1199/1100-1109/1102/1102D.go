package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	bytes := []byte(s)
	target := n / 3
	a, b, c := 0, 0, 0
	for _, ch := range bytes {
		switch ch {
		case '0':
			a++
		case '1':
			b++
		case '2':
			c++
		}
	}
	// increase zeros
	for i := 0; i < n && a < target; i++ {
		if bytes[i] == '1' && b > target {
			bytes[i], a, b = '0', a+1, b-1
		} else if bytes[i] == '2' && c > target {
			bytes[i], a, c = '0', a+1, c-1
		}
	}
	// mark extra zeros
	for i := n - 1; i >= 0 && a > target; i-- {
		if bytes[i] == '0' {
			bytes[i] = '3'
			a--
		}
	}
	// fill ones from placeholders
	for i := 0; i < n && b < target; i++ {
		if bytes[i] == '3' {
			bytes[i] = '1'
			b++
		}
	}
	// fill ones from twos
	for i := 0; i < n && b < target; i++ {
		if bytes[i] == '2' {
			bytes[i], b, c = '1', b+1, c-1
		}
	}
	// mark extra ones
	for i := n - 1; i >= 0 && b > target; i-- {
		if bytes[i] == '1' {
			bytes[i] = '3'
			b--
		}
	}
	// final replace placeholders with '2'
	for i := 0; i < n; i++ {
		if bytes[i] == '3' {
			bytes[i] = '2'
		}
	}
	// output result
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	writer.Write(bytes)
}
