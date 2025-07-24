package main

import (
	"bufio"
	"fmt"
	"os"
)

func simulate(orig []byte, start int) int {
	n := len(orig)
	s := make([]byte, n)
	copy(s, orig)
	pos := start
	steps := 0
	limit := n*n*2 + 5
	for steps <= limit {
		if pos < 1 || pos > n {
			return steps
		}
		if s[pos-1] == 'U' {
			s[pos-1] = 'D'
			pos++
		} else {
			s[pos-1] = 'U'
			pos--
		}
		steps++
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var str string
	fmt.Fscan(reader, &str)
	arr := []byte(str)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		res := simulate(arr, i)
		fmt.Fprint(writer, res)
	}
	fmt.Fprintln(writer)
}
