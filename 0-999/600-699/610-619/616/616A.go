package main

import (
	"bufio"
	"fmt"
	"os"
)

func trimLeadingZeros(s string) string {
	i := 0
	for i < len(s) && s[i] == '0' {
		i++
	}
	if i == len(s) {
		return "0"
	}
	return s[i:]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b string
	fmt.Fscan(reader, &a)
	fmt.Fscan(reader, &b)
	a = trimLeadingZeros(a)
	b = trimLeadingZeros(b)
	if len(a) < len(b) {
		fmt.Println("<")
	} else if len(a) > len(b) {
		fmt.Println(">")
	} else {
		if a == b {
			fmt.Println("=")
		} else if a < b {
			fmt.Println("<")
		} else {
			fmt.Println(">")
		}
	}
}
