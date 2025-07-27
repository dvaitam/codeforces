package main

import (
	"bufio"
	"fmt"
	"os"
)

func isValid(candidate []byte, arr [][]byte, m int) bool {
	for _, s := range arr {
		diff := 0
		for j := 0; j < m && diff <= 1; j++ {
			if candidate[j] != s[j] {
				diff++
			}
		}
		if diff > 1 {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		arr := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			arr[i] = []byte(s)
		}
		found := false
		result := ""
		for i := 0; i < n && !found; i++ {
			base := make([]byte, m)
			copy(base, arr[i])
			if isValid(base, arr, m) {
				found = true
				result = string(base)
				break
			}
			for j := 0; j < m && !found; j++ {
				orig := base[j]
				for c := byte('a'); c <= 'z'; c++ {
					base[j] = c
					if isValid(base, arr, m) {
						found = true
						result = string(base)
						break
					}
				}
				base[j] = orig
			}
		}
		if found {
			fmt.Fprintln(writer, result)
		} else {
			fmt.Fprintln(writer, "-1")
		}
	}
}
