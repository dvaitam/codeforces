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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		var s string
		fmt.Fscan(reader, &s)
		n := len(s)
		seen := make([]bool, 26)
		minIdx, maxIdx := 25, 0
		ok := true
		for _, ch := range s {
			idx := int(ch - 'a')
			if idx < 0 || idx >= 26 || seen[idx] {
				ok = false
				break
			}
			seen[idx] = true
			if idx < minIdx {
				minIdx = idx
			}
			if idx > maxIdx {
				maxIdx = idx
			}
		}
		if !ok || maxIdx-minIdx+1 != n {
			fmt.Fprintln(writer, "No")
		} else {
			fmt.Fprintln(writer, "Yes")
		}
	}
}
