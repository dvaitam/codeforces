package main

import (
	"bufio"
	"fmt"
	"os"
)

func steps(s []byte, start int) int {
	n := len(s)
	// create copy to avoid mutating original string
	b := make([]byte, n)
	copy(b, s)
	pos := start
	t := 0
	for pos >= 0 && pos < n {
		if b[pos] == '>' {
			b[pos] = '<'
			pos++
		} else {
			b[pos] = '>'
			pos--
		}
		t++
	}
	return t
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	var str string
	fmt.Fscan(in, &str)
	s := []byte(str)
	for i := 0; i < n; i++ {
		fmt.Fprint(out, steps(s, i))
		if i+1 < n {
			fmt.Fprint(out, " ")
		}
	}
	fmt.Fprintln(out)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve(in, out)
	}
}
