package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	for i := 0; i+1 < n; i++ {
		if s[i] == '1' && s[i+1] == '1' {
			fmt.Println("No")
			return
		}
	}
	for i := 0; i < n; i++ {
		if s[i] == '0' {
			leftEmpty := i == 0 || s[i-1] == '0'
			rightEmpty := i == n-1 || s[i+1] == '0'
			if leftEmpty && rightEmpty {
				fmt.Println("No")
				return
			}
		}
	}
	fmt.Println("Yes")
}
