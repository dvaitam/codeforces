package main

import (
	"bufio"
	"fmt"
	"os"
)

func isValid(s []byte) bool {
	bal := 0
	for _, c := range s {
		if c == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		openCnt, closeCnt := 0, 0
		for i := 0; i < n; i++ {
			if s[i] == '(' {
				openCnt++
			} else if s[i] == ')' {
				closeCnt++
			}
		}
		openNeed := n/2 - openCnt
		closeNeed := n/2 - closeCnt
		bytes := []byte(s)
		openPos := []int{}
		closePos := []int{}
		for i := 0; i < n; i++ {
			if bytes[i] == '?' {
				if openNeed > 0 {
					bytes[i] = '('
					openNeed--
					openPos = append(openPos, i)
				} else {
					bytes[i] = ')'
					closeNeed--
					closePos = append(closePos, i)
				}
			}
		}
		if len(openPos) == 0 || len(closePos) == 0 {
			if isValid(bytes) {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
			continue
		}
		// swap last '(' we inserted with first ')' we inserted
		i := openPos[len(openPos)-1]
		j := closePos[0]
		bytes[i] = ')'
		bytes[j] = '('
		if isValid(bytes) {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
