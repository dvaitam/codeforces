package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s, t string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &s)
	fmt.Fscan(in, &t)

	sb := []byte(s)
	tb := []byte(t)

	l := 0
	for l < n && sb[l] == tb[l] {
		l++
	}
	if l == n {
		fmt.Println(0)
		return
	}

	r := n - 1
	for r >= 0 && sb[r] == tb[r] {
		r--
	}

	ans := make(map[string]struct{}, 2)

	if checkShift(sb, tb, l, r) {
		w := buildInsert(s, l, tb[l])
		ans[w] = struct{}{}
	}

	if checkShift(tb, sb, l, r) {
		w := buildInsert(t, l, sb[l])
		ans[w] = struct{}{}
	}

	fmt.Println(len(ans))
}

func checkShift(a, b []byte, l, r int) bool {
	for i := l; i < r; i++ {
		if a[i] != b[i+1] {
			return false
		}
	}
	return true
}

func buildInsert(base string, pos int, ch byte) string {
	var b strings.Builder
	b.Grow(len(base) + 1)
	b.WriteString(base[:pos])
	b.WriteByte(ch)
	b.WriteString(base[pos:])
	return b.String()
}
