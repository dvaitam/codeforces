package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var phase string
	if _, err := fmt.Fscan(in, &phase); err != nil {
		return
	}

	if phase == "first" {
		handleFirstRun(in, out)
	} else if phase == "second" {
		handleSecondRun(in, out)
	}
}

func handleFirstRun(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	header := make([]byte, 3)
	value := n
	for i := 2; i >= 0; i-- {
		header[i] = byte('a' + (value % 26))
		value /= 26
	}

	res := make([]byte, 3+n)
	copy(res, header)
	for i, v := range arr {
		res[3+i] = byte('a' + v - 1)
	}

	fmt.Fprintln(out, string(res))
}

func handleSecondRun(in *bufio.Reader, out *bufio.Writer) {
	var s string
	fmt.Fscan(in, &s)
	if len(s) < 3 {
		fmt.Fprintln(out, 0)
		fmt.Fprintln(out)
		return
	}

	n := 0
	for i := 0; i < 3; i++ {
		n = n*26 + int(s[i]-'a')
	}

	arr := make([]int, len(s)-3)
	for i := 3; i < len(s); i++ {
		arr[i-3] = int(s[i]-'a') + 1
	}

	fmt.Fprintln(out, n)
	for i, v := range arr {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
