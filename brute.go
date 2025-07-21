package main
import (
	"fmt"
	"math"
)
func main() {
	maxK := 12
	for K := 1; K <= maxK; K++ {
		n := K
		A := n + 1
		// coordinates function
		type P struct{ x, y float64 }
		var pts1, pts2, pts3 []P
		alpha := math.Sqrt(3) / 2.0
		for i := 1; i <= n; i++ {
			x := float64(i) / float64(A)
			pts1 = append(pts1, P{ x, 0 })
			// side 2
			y := float64(i) / float64(A)
			pts2 = append(pts2, P{ 1 - 0.5*y, alpha * y })
			// side 3
			z := float64(i) / float64(A)
			pts3 = append(pts3, P{ 0.5*(1 - z), alpha * (1 - z) })
		}
		cnt := 0
		for _, p1 := range pts1 {
			for _, p2 := range pts2 {
				for _, p3 := range pts3 {
					// check obtuse
					d12 := (p2.x-p1.x)*(p3.x-p1.x) + (p2.y-p1.y)*(p3.y-p1.y)
					if d12 < 0 {
						cnt++
					} else {
						// check at p2
						d21 := (p1.x-p2.x)*(p3.x-p2.x) + (p1.y-p2.y)*(p3.y-p2.y)
						if d21 < 0 {
							cnt++
						} else {
							// at p3
							d31 := (p1.x-p3.x)*(p2.x-p3.x) + (p1.y-p3.y)*(p2.y-p3.y)
							if d31 < 0 {
								cnt++
							}
						}
					}
				}
			}
		}
		fmt.Printf("%d %d\n", K, cnt)
	}
}
