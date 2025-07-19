package main

import "fmt"

func main() {
	var n, l, k int
	if _, err := fmt.Scan(&n, &l, &k); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&p[i])
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&c[i])
	}
	var b, b_ [201]float64
	b1 := &b
	b2 := &b_
	(*b1)[0] = 1.0
	kb := 0
	for i := 0; i < n; i++ {
		if c[i] == -1 {
			p1 := float64(p[i]) / 100.0
			p0 := float64(100-p[i]) / 100.0
			for j := 0; j <= kb+1; j++ {
				(*b2)[j] = 0
			}
			for j := 0; j <= kb; j++ {
				(*b2)[j] += p0 * (*b1)[j]
				(*b2)[j+1] += p1 * (*b1)[j]
			}
			b1, b2 = b2, b1
			kb++
		}
	}
	ma := kb - k
	if ma < 0 {
		ma = 0
	}
	var a, a_ [201][201]float64
	a1 := &a
	a2 := &a_
	(*a1)[0][0] = 1.0
	ka := 0
	for i := 0; i < n; i++ {
		if c[i] != -1 {
			p1 := float64(p[i]) / 100.0
			p0 := float64(100-p[i]) / 100.0
			for x := 0; x <= ka+1; x++ {
				for y := 0; y <= ma; y++ {
					(*a2)[x][y] = 0
				}
			}
			for x := 0; x <= ka; x++ {
				for y := 0; y <= ma; y++ {
					(*a2)[x][y] += p0 * (*a1)[x][y]
					idx := y + c[i]
					if idx > ma {
						idx = ma
					}
					(*a2)[x+1][idx] += p1 * (*a1)[x][y]
				}
			}
			a1, a2 = a2, a1
			ka++
		}
	}
	ans := 0.0
	for i := 0; i <= ka; i++ {
		for j := 0; j <= ma; j++ {
			for t := 0; t <= kb; t++ {
				if t+i >= l && j+k >= t {
					ans += (*a1)[i][j] * (*b1)[t]
				}
			}
		}
	}
	fmt.Printf("%.20f\n", ans)
}
