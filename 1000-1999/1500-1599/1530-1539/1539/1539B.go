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

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	pref := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + int(s[i-1]-'a'+1)
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		ans := pref[r] - pref[l-1]
		fmt.Fprintln(writer, ans)
	}
}
