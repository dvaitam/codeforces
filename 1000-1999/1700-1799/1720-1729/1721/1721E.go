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

	var s string
	fmt.Fscan(reader, &s)
	var q int
	fmt.Fscan(reader, &q)

	sb := []byte(s)
	n := len(sb)

	// prefix function for s
	piS := make([]int, n)
	for i := 1; i < n; i++ {
		j := piS[i-1]
		for j > 0 && sb[i] != sb[j] {
			j = piS[j-1]
		}
		if sb[i] == sb[j] {
			j++
		}
		piS[i] = j
	}

	for ; q > 0; q-- {
		var t string
		fmt.Fscan(reader, &t)
		tb := []byte(t)
		piT := make([]int, len(tb))
		j := piS[n-1]
		for i := 0; i < len(tb); i++ {
			c := tb[i]
			get := func(pos int) byte {
				if pos < n {
					return sb[pos]
				}
				return tb[pos-n]
			}
			for j > 0 && get(j) != c {
				if j <= n {
					j = piS[j-1]
				} else {
					j = piT[j-n-1]
				}
			}
			if get(j) == c {
				j++
			}
			piT[i] = j
		}
		for i := 0; i < len(piT); i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, piT[i])
		}
		writer.WriteByte('\n')
	}
}
