package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func prefixRepeat(s string, d int) int {
	var b strings.Builder
	for b.Len() < d {
		b.WriteString(s)
	}
	str := b.String()[:d]
	val, _ := strconv.Atoi(str)
	return val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ns := strconv.Itoa(n)
		l := len(ns)
		type pair struct{ a, b int }
		res := make([]pair, 0)
		for a := 1; a <= 10000; a++ {
			for d := 1; d <= 7; d++ {
				b := l*a - d
				if b < 1 || b > a*n || b > 10000 {
					continue
				}
				if d > l*a {
					continue
				}
				pref := prefixRepeat(ns, d)
				if pref == a*n-b {
					res = append(res, pair{a, b})
				}
			}
		}
		fmt.Fprintln(out, len(res))
		for _, p := range res {
			fmt.Fprintln(out, p.a, p.b)
		}
	}
}
