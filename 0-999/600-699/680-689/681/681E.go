package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type segment struct{ l, r float64 }

func main() {
	in := bufio.NewReader(os.Stdin)
	var x0, y0, v, T int64
	if _, err := fmt.Fscan(in, &x0, &y0, &v, &T); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	R := float64(v) * float64(T)
	if R > 6e9 {
		R = 6e9
	}
	segs := make([]segment, 0, n)
	inside := false
	for i := 0; i < n; i++ {
		var xi, yi, ri int64
		fmt.Fscan(in, &xi, &yi, &ri)
		dx := float64(xi - x0)
		dy := float64(yi - y0)
		dist2 := dx*dx + dy*dy
		r := float64(ri)
		if dist2 <= r*r {
			inside = true
			break
		}
		dist := math.Sqrt(dist2)
		if r == 0 || R+r <= dist {
			continue
		}
		angl := math.Atan2(dx, dy)
		al := math.Asin(r / dist)
		if R*R+r*r < dist*dist {
			cosv := (R*R + dist*dist - r*r) / (2 * R * dist)
			if cosv > 1 {
				cosv = 1
			} else if cosv < -1 {
				cosv = -1
			}
			al = math.Acos(cosv)
		}
		ang1 := angl - al
		ang2 := angl + al
		if ang1 < -math.Pi {
			segs = append(segs, segment{-math.Pi, ang2})
			segs = append(segs, segment{ang1 + 2*math.Pi, math.Pi})
		} else if ang2 > math.Pi {
			segs = append(segs, segment{ang1, math.Pi})
			segs = append(segs, segment{-math.Pi, ang2 - 2*math.Pi})
		} else {
			segs = append(segs, segment{ang1, ang2})
		}
	}
	if inside {
		fmt.Printf("%.9f\n", 1.0)
		return
	}
	if len(segs) == 0 {
		fmt.Printf("%.9f\n", 0.0)
		return
	}
   sort.Slice(segs, func(i, j int) bool {
       if segs[i].l != segs[j].l {
           return segs[i].l < segs[j].l
       }
       return segs[i].r < segs[j].r
   })
	var total, currLen, end float64
	end = -math.Pi
	for _, s := range segs {
		if end >= s.l {
			if end < s.r {
				currLen += s.r - end
				end = s.r
			}
		} else {
			total += currLen
			currLen = s.r - s.l
			end = s.r
		}
	}
	total += currLen
	p := total / (2 * math.Pi)
	fmt.Printf("%.9f\n", p)
}
