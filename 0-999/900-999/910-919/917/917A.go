package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
		if len(s) > 0 && s[len(s)-1] == '\r' {
			s = s[:len(s)-1]
		}
	}
	n := len(s)
	ans := 0
	for i := 0; i < n; i++ {
		minBal := 0
		maxBal := 0
		for j := i; j < n; j++ {
			switch s[j] {
			case '(':
				minBal++
				maxBal++
			case ')':
				minBal--
				maxBal--
			default:
				minBal--
				maxBal++
			}
			if maxBal < 0 {
				break
			}
			if minBal < 0 {
				minBal = 0
			}
			if (j-i+1)%2 == 0 && minBal == 0 {
				ans++
			}
		}
	}
	fmt.Println(ans)
}
