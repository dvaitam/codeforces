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

	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	if len(s) == 0 {
		fmt.Fprintln(writer, 0)
		return
	}
	// count transitions and equal adjacent pairs
	transitions := 0
	equalCnt := 0
	for i := 0; i+1 < len(s); i++ {
		if s[i] != s[i+1] {
			transitions++
		} else {
			equalCnt++
		}
	}

	// after flipping one substring, at most two boundaries change
	delta := 0
	if equalCnt >= 2 {
		delta = 2
	} else if equalCnt == 1 {
		delta = 1
	}

	result := transitions + delta + 1
	fmt.Fprintln(writer, result)
}
