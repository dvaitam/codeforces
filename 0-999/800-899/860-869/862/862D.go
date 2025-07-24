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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if len(s) < n {
		// If the string may contain spaces, read the entire line
		rest, _ := in.ReadString('\n')
		s += rest
		if len(s) > n {
			s = s[:n]
		}
	}

	pos0, pos1 := -1, -1
	for i, ch := range s {
		if ch == '0' && pos0 == -1 {
			pos0 = i + 1
		}
		if ch == '1' && pos1 == -1 {
			pos1 = i + 1
		}
	}

	fmt.Fprintln(out, pos0, pos1)
}
