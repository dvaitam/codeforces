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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		cnt := map[byte]int{'L': 0, 'I': 0, 'T': 0}
		for i := 0; i < n; i++ {
			cnt[s[i]]++
		}
		if cnt['L'] == cnt['I'] && cnt['L'] == cnt['T'] {
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, -1)
	}
}
