package main

import (
	"bufio"
	"fmt"
	"os"
)

func isRegular(s []byte) bool {
	bal := 0
	for _, c := range s {
		if c == '(' {
			bal++
		} else {
			bal--
			if bal < 0 {
				return false
			}
		}
	}
	return bal == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	arr := []byte(s)
	cnt := 0
	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			for i := l; i <= r; i++ {
				if arr[i] == '(' {
					arr[i] = ')'
				} else {
					arr[i] = '('
				}
			}
			if isRegular(arr) {
				cnt++
			}
			for i := l; i <= r; i++ {
				if arr[i] == '(' {
					arr[i] = ')'
				} else {
					arr[i] = '('
				}
			}
		}
	}
	fmt.Fprintln(writer, cnt)
}
