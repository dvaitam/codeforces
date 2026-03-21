package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func getCost(x, y, D, attacks int64) int64 {
	if x > 0 && D > 2000000000000000000/x {
		return 2000000000000000000
	}
	return x*D + y*attacks
}

func solve() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var x, y, z, k int64
		fmt.Fscan(reader, &x, &y, &z, &k)

		ans := int64(2000000000000000000)

		for q := int64(0); ; q++ {
			Rq := z - q*(q+1)/2*k
			if Rq <= 0 {
				low := int64(1)
				high := q * k
				if high < 1 {
					high = 1
				}
				ans_D := q * k
				for low <= high {
					mid := low + (high-low)/2
					c := mid / k
					if c > q {
						c = q
					}
					dmg := c*(c+1)/2*k + (q-c)*mid
					if dmg >= z {
						ans_D = mid
						high = mid - 1
					} else {
						low = mid + 1
					}
				}
				cost := getCost(x, y, ans_D, q)
				if cost < ans {
					ans = cost
				}
				break
			}

			L := q * k
			if L < 1 {
				L = 1
			}
			R_max := (q+1)*k - 1

			pL := (Rq + L - 1) / L
			costL := getCost(x, y, L, q+pL)
			if costL < ans {
				ans = costL
			}

			pR := (Rq + R_max - 1) / R_max
			costR := getCost(x, y, R_max, q+pR)
			if costR < ans {
				ans = costR
			}

			if x >= y {
				D_unc := math.Sqrt(float64(y) * float64(Rq) / float64(x))
				Dc := int64(D_unc)
				D_center := Dc
				if D_center < L {
					D_center = L
				}
				if D_center > R_max {
					D_center = R_max
				}
				start := D_center - 500
				if start < L {
					start = L
				}
				end := D_center + 500
				if end > R_max {
					end = R_max
				}
				for D := start; D <= end; D++ {
					p := (Rq + D - 1) / D
					cost := getCost(x, y, D, q+p)
					if cost < ans {
						ans = cost
					}
				}
			} else {
				p_unc := math.Sqrt(float64(x) * float64(Rq) / float64(y))
				pc := int64(p_unc)
				P_min := (Rq + R_max - 1) / R_max
				if P_min < 1 {
					P_min = 1
				}
				P_max := (Rq + L - 1) / L

				p_center := pc
				if p_center < P_min {
					p_center = P_min
				}
				if p_center > P_max {
					p_center = P_max
				}
				start := p_center - 500
				if start < 1 {
					start = 1
				}
				end := p_center + 500
				for p := start; p <= end; p++ {
					D := (Rq + p - 1) / p
					if D < L {
						D = L
					}
					if D <= R_max {
						cost := getCost(x, y, D, q+p)
						if cost < ans {
							ans = cost
						}
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}

func main() {
	solve()
}
