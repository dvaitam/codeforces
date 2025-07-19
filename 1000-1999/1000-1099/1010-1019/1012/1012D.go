package main

import (
	"bufio"
	"fmt"
	"os"
)

// Pair holds two integer values
type Pair struct{ first, second int }

// reverseInts reverses a slice of ints in place
func reverseInts(a []int) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// get computes the operations for strings s and t as in solD.cpp logic
func get(s, t string) []Pair {
	// trim trailing 'a' from s and trailing 'b' from t
	for len(s) > 0 && s[len(s)-1] == 'a' {
		s = s[:len(s)-1]
	}
	for len(t) > 0 && t[len(t)-1] == 'b' {
		t = t[:len(t)-1]
	}
	// run-length encode s and t
	sv, tv := []int{}, []int{}
	cnt := 0
	for i := 0; i < len(s); i++ {
		cnt++
		if i == len(s)-1 || s[i+1] != s[i] {
			sv = append(sv, cnt)
			cnt = 0
		}
	}
	cnt = 0
	for i := 0; i < len(t); i++ {
		cnt++
		if i == len(t)-1 || t[i+1] != t[i] {
			tv = append(tv, cnt)
			cnt = 0
		}
	}
	reverseInts(sv)
	reverseInts(tv)
	if len(sv) == 0 && len(tv) == 0 {
		return nil
	}
	fl := false
	if len(tv) == 0 {
		sv, tv = tv, sv
		fl = true
	}
	ans := []Pair{}
	// handle empty sv
	if len(sv) == 0 {
		m := len(tv) % 4
		if m == 0 || m == 1 {
			sum := tv[len(tv)/2]
			for i := len(tv)/2 + 1; i < len(tv); i++ {
				sum += tv[i]
				sv = append(sv, tv[i])
			}
			tv = tv[:len(tv)/2]
			ans = append(ans, Pair{0, sum})
		} else {
			sum := 0
			for i := len(tv) / 2; i < len(tv); i++ {
				sum += tv[i]
				sv = append(sv, tv[i])
			}
			tv = tv[:len(tv)/2]
			ans = append(ans, Pair{0, sum})
		}
	}
	// main reduction loop
	for len(sv) > 0 && len(tv) > 0 {
		switch {
		case len(sv) == 1:
			if len(tv) == 1 {
				ans = append(ans, Pair{sv[0], tv[0]})
				sv = sv[:0]
				tv = tv[:0]
				continue
			}
			if len(tv)%2 == 0 {
				a := sv[len(sv)-1]
				b := tv[len(tv)-1] + tv[len(tv)-2]
				ans = append(ans, Pair{a, b})
				if len(tv) > 2 {
					tv[len(tv)-3] += a
				}
				last := tv[len(tv)-1]
				sv = sv[:len(sv)-1]
				sv = append(sv, last)
				tv = tv[:len(tv)-2]
			} else {
				a := sv[len(sv)-1]
				b := tv[len(tv)-1] + tv[len(tv)-2] + tv[len(tv)-3]
				ans = append(ans, Pair{a, b})
				if len(tv) > 3 {
					tv[len(tv)-4] += a
				}
				t1, t2 := tv[len(tv)-2], tv[len(tv)-1]
				sv = sv[:len(sv)-1]
				sv = append(sv, t1, t2)
				tv = tv[:len(tv)-3]
			}
		case len(tv) == 1:
			if len(sv)%2 == 0 {
				a := sv[len(sv)-1] + sv[len(sv)-2]
				b := tv[0]
				ans = append(ans, Pair{a, b})
				if len(sv) > 2 {
					sv[len(sv)-3] += b
				}
				last := sv[len(sv)-1]
				sv = sv[:len(sv)-2]
				tv = []int{last}
			} else {
				a := sv[len(sv)-1] + sv[len(sv)-2] + sv[len(sv)-3]
				b := tv[0]
				ans = append(ans, Pair{a, b})
				if len(sv) > 3 {
					sv[len(sv)-4] += b
				}
				b2, b1 := sv[len(sv)-2], sv[len(sv)-1]
				sv = sv[:len(sv)-3]
				tv = []int{b2, b1}
			}
		case len(sv)%2 == len(tv)%2:
			a := sv[len(sv)-1]
			b := tv[len(tv)-1]
			ans = append(ans, Pair{a, b})
			sv[len(sv)-2] += b
			tv[len(tv)-2] += a
			sv = sv[:len(sv)-1]
			tv = tv[:len(tv)-1]
		default:
			if len(sv) <= len(tv) {
				a := sv[len(sv)-1]
				b := tv[len(tv)-1] + tv[len(tv)-2]
				ans = append(ans, Pair{a, b})
				sv[len(sv)-2] += tv[len(tv)-2]
				tv[len(tv)-3] += a
				last := tv[len(tv)-1]
				sv = sv[:len(sv)-1]
				sv = append(sv, last)
				tv = tv[:len(tv)-2]
			} else {
				a := sv[len(sv)-1] + sv[len(sv)-2]
				b := tv[len(tv)-1]
				ans = append(ans, Pair{a, b})
				tv[len(tv)-2] += sv[len(sv)-2]
				sv[len(sv)-3] += b
				last := sv[len(sv)-1]
				sv = sv[:len(sv)-2]
				tv = tv[:len(tv)-1]
				tv = append(tv, last)
			}
		}
	}
	// handle leftovers
	if len(sv) == 0 {
		if len(tv) == 1 {
			ans = append(ans, Pair{0, tv[0]})
		} else if len(tv) == 2 {
			ans = append(ans, Pair{0, tv[1] + tv[0]})
			ans = append(ans, Pair{tv[1], 0})
		}
	} else {
		if len(sv) == 1 {
			ans = append(ans, Pair{sv[0], 0})
		} else if len(sv) == 2 {
			ans = append(ans, Pair{sv[1] + sv[0], 0})
			ans = append(ans, Pair{0, sv[1]})
		}
	}
	// if swapped, swap pairs
	if fl {
		for i := range ans {
			ans[i] = Pair{ans[i].second, ans[i].first}
		}
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var s, t string
	fmt.Fscan(reader, &s, &t)
	vv := get(s, t)
	vv2 := get(t, s)
	if len(vv2) < len(vv) {
		vv = vv2
		for i := range vv {
			vv[i] = Pair{vv[i].second, vv[i].first}
		}
	}
	fmt.Fprintln(writer, len(vv))
	for _, e := range vv {
		fmt.Fprintln(writer, e.first, e.second)
	}
}
