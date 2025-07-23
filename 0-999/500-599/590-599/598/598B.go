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

	var m int
	fmt.Fscan(in, &m)

	b := []byte(s)
	for i := 0; i < m; i++ {
		var l, r, k int
		fmt.Fscan(in, &l, &r, &k)
		l--
		length := r - l
		k %= length
		if k > 0 {
			tmp := append([]byte(nil), b[l:r]...)
			copy(b[l:r], append(tmp[length-k:], tmp[:length-k]...))
		}
	}

	fmt.Fprintln(out, string(b))
}
