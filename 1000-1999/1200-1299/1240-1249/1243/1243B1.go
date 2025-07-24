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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)

		mism := make([]int, 0, 2)
		for i := 0; i < n; i++ {
			if s[i] != t[i] {
				mism = append(mism, i)
			}
		}

		if len(mism) == 2 && s[mism[0]] == s[mism[1]] && t[mism[0]] == t[mism[1]] {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
