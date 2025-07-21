package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	s := strings.TrimSpace(line)

	// find operator position
	var op rune
	var pos int
	for i, ch := range s {
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' {
			op = ch
			pos = i
			break
		}
	}

	a, _ := strconv.Atoi(s[:pos])
	b, _ := strconv.Atoi(s[pos+1:])

	var res int
	switch op {
	case '+':
		res = a + b
	case '-':
		res = a - b
	case '*':
		res = a * b
	case '/':
		res = a / b
	case '%':
		res = a % b
	}

	fmt.Println(res)
}
