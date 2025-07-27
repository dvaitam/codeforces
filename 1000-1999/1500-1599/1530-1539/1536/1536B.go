package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	substr1 := make(map[string]struct{})
	substr2 := make(map[string]struct{})
	substr3 := make(map[string]struct{})

	for i := 0; i < n; i++ {
		substr1[s[i:i+1]] = struct{}{}
		if i+1 < n {
			substr2[s[i:i+2]] = struct{}{}
		}
		if i+2 < n {
			substr3[s[i:i+3]] = struct{}{}
		}
	}

	for c := byte('a'); c <= 'z'; c++ {
		str := string([]byte{c})
		if _, ok := substr1[str]; !ok {
			fmt.Fprintln(writer, str)
			return
		}
	}

	for c1 := byte('a'); c1 <= 'z'; c1++ {
		for c2 := byte('a'); c2 <= 'z'; c2++ {
			str := string([]byte{c1, c2})
			if _, ok := substr2[str]; !ok {
				fmt.Fprintln(writer, str)
				return
			}
		}
	}

	for c1 := byte('a'); c1 <= 'z'; c1++ {
		for c2 := byte('a'); c2 <= 'z'; c2++ {
			for c3 := byte('a'); c3 <= 'z'; c3++ {
				str := string([]byte{c1, c2, c3})
				if _, ok := substr3[str]; !ok {
					fmt.Fprintln(writer, str)
					return
				}
			}
		}
	}
}
