package main

import (
	"bufio"
	"fmt"
	"os"
)

type Var struct {
	pref string
	suff string
	cnt  int
}

func countHaha(s string) int {
	c := 0
	for i := 0; i+3 < len(s); i++ {
		if s[i:i+4] == "haha" {
			c++
		}
	}
	return c
}

func merge(a, b Var) Var {
	res := Var{}
	res.cnt = a.cnt + b.cnt + countHaha(a.suff+b.pref)
	pref := a.pref + b.pref
	if len(pref) > 3 {
		pref = pref[:3]
	}
	res.pref = pref
	suff := a.suff + b.suff
	if len(suff) > 3 {
		suff = suff[len(suff)-3:]
	}
	res.suff = suff
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		vars := make(map[string]Var)
		last := ""
		for i := 0; i < n; i++ {
			var x, op string
			fmt.Fscan(in, &x, &op)
			if op == ":=" {
				var s string
				fmt.Fscan(in, &s)
				v := Var{pref: s, suff: s, cnt: countHaha(s)}
				if len(v.pref) > 3 {
					v.pref = v.pref[:3]
				}
				if len(v.suff) > 3 {
					v.suff = v.suff[len(v.suff)-3:]
				}
				vars[x] = v
			} else { // op == "="
				var a, plus, b string
				fmt.Fscan(in, &a, &plus, &b)
				v := merge(vars[a], vars[b])
				vars[x] = v
			}
			last = x
		}
		fmt.Fprintln(out, vars[last].cnt)
	}
}
