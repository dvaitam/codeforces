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
		prefix := []byte{s[0]}
		for i := 1; i < n && s[i] <= s[i-1]; i++ {
			prefix = append(prefix, s[i])
		}
		if len(prefix) > 1 && prefix[0] == prefix[1] {
			prefix = prefix[:1]
		}
		res := make([]byte, len(prefix)*2)
		copy(res, prefix)
		for i := 0; i < len(prefix); i++ {
			res[len(prefix)+i] = prefix[len(prefix)-1-i]
		}
		fmt.Fprintln(out, string(res))
	}
}
