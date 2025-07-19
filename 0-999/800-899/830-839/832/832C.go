package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const Maxn = 1000000

type node struct {
	x, v int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, s int
	for {
		if _, err := fmt.Fscan(in, &n, &s); err != nil {
			break
		}
		a1 := make([]node, 0, n)
		a2 := make([]node, 0, n)
		for i := 0; i < n; i++ {
			var x, v, t int
			fmt.Fscan(in, &x, &v, &t)
			if t == 1 {
				a1 = append(a1, node{x, v})
			} else {
				a2 = append(a2, node{x, v})
			}
		}
		sort.Slice(a1, func(i, j int) bool {
			if a1[i].x == a1[j].x {
				return a1[i].v > a1[j].v
			}
			return a1[i].x < a1[j].x
		})
		sort.Slice(a2, func(i, j int) bool {
			if a2[i].x == a2[j].x {
				return a2[i].v > a2[j].v
			}
			return a2[i].x < a2[j].x
		})

		s1 := make([]int, Maxn+1)
		s2 := make([]int, Maxn+1)
		// process a1
		pos, sm := 0, -1
		vm := 1e7
		t1 := 1e7
		for i := 0; i < len(a1); i++ {
			for pos < a1[i].x {
				s1[pos] = sm
				pos++
			}
			val1 := 1.0 / float64(a1[i].v+s) * float64(a1[i].x)
			if vm > val1 {
				vm = val1
				sm = i
			}
			ttmp := 1.0 / float64(a1[i].v) * float64(a1[i].x)
			if ttmp < t1 {
				t1 = ttmp
			}
		}
		for pos <= Maxn {
			s1[pos] = sm
			pos++
		}
		// process a2
		pos = Maxn
		sm = -1
		vm = 1e7
		t2 := 1e7
		for i := len(a2) - 1; i >= 0; i-- {
			for pos > a2[i].x {
				s2[pos] = sm
				pos--
			}
			val2 := 1.0 / float64(a2[i].v+s) * float64(a2[i].x)
			if vm > val2 {
				vm = val2
				sm = i
			}
			ttmp2 := 1.0 / float64(a2[i].v) * float64(Maxn-a2[i].x)
			if ttmp2 < t2 {
				t2 = ttmp2
			}
		}
		for pos >= 0 {
			s2[pos] = sm
			pos--
		}
		// combine
		ans := 1e7
		for i := 0; i <= Maxn; i++ {
			tt1, tt2 := t1, t2
			if s1[i] >= 0 {
				idx := s1[i]
				v := a1[idx].v
				x := a1[idx].x
				coeff := 1.0 / float64(s+v) * 1.0 / float64(s-v)
				val := coeff * (float64(i*s) - float64(v*x))
				if val < tt1 {
					tt1 = val
				}
			}
			if s2[i] >= 0 {
				idx := s2[i]
				v := a2[idx].v
				x := a2[idx].x
				coeff := 1.0 / float64(s+v) * 1.0 / float64(s-v)
				val := coeff * (float64((Maxn-i)*s) - float64(v*(Maxn-x)))
				if val < tt2 {
					tt2 = val
				}
			}
			// take max of two
			var m float64
			if tt1 > tt2 {
				m = tt1
			} else {
				m = tt2
			}
			if m < ans {
				ans = m
			}
		}
		fmt.Printf("%.12f\n", ans)
	}
}
