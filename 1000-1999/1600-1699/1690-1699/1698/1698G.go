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

	var s string
	fmt.Fscan(in, &s)

	idx := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			idx = i + 1 // 1-indexed position
			break
		}
	}
	if idx == -1 {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, idx, idx+1)
}
