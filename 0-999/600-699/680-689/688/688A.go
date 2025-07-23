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

	var n, d int
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}
	maxSeq := 0
	curSeq := 0
	for i := 0; i < d; i++ {
		var s string
		fmt.Fscan(in, &s)
		ok := false
		for j := 0; j < len(s); j++ {
			if s[j] == '0' {
				ok = true
				break
			}
		}
		if ok {
			curSeq++
			if curSeq > maxSeq {
				maxSeq = curSeq
			}
		} else {
			curSeq = 0
		}
	}
	fmt.Fprintln(out, maxSeq)
}
