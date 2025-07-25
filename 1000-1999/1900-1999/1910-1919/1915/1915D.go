package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

var (
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		res := splitSyllables(s)
		fmt.Fprintln(out, res)
	}
}

func splitSyllables(s string) string {
	isVowel := func(c byte) bool {
		return c == 'a' || c == 'e'
	}
	n := len(s)
	vPos := []int{}
	for i := 0; i < n; i++ {
		if isVowel(s[i]) {
			vPos = append(vPos, i)
		}
	}
	boundaries := make(map[int]bool)
	for i := 0; i+1 < len(vPos); i++ {
		diff := vPos[i+1] - vPos[i]
		if diff == 2 {
			boundaries[vPos[i]] = true
		} else if diff == 3 {
			boundaries[vPos[i]+1] = true
		}
	}
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteByte(s[i])
		if boundaries[i] {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}
