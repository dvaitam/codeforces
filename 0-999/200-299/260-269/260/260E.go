package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	x, y int
}

type YRange struct {
	y1min, y1max float64
	y2min, y2max float64
}

func main() {
	// Fast I/O
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	n := nextInt()
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		points[i].x = nextInt()
		points[i].y = nextInt()
	}

	a := make([]int, 9)
	for i := 0; i < 9; i++ {
		a[i] = nextInt()
	}

	// Sort points by X
	pointsX := make([]Point, n)
	copy(pointsX, points)
	sort.Slice(pointsX, func(i, j int) bool {
		return pointsX[i].x < pointsX[j].x
	})

	// Sort points by Y
	pointsY := make([]Point, n)
	copy(pointsY, points)
	sort.Slice(pointsY, func(i, j int) bool {
		return pointsY[i].y < pointsY[j].y
	})

	// Pre-allocate buffers for partitions
	p1 := make([]Point, 0, n)
	p2 := make([]Point, 0, n)
	p3 := make([]Point, 0, n)

	// Generate bitmasks for combinations of 3 items from 9
	masks3 := make([]int, 0, 84)
	for i := 0; i < (1 << 9); i++ {
		if popCount(i) == 3 {
			masks3 = append(masks3, i)
		}
	}

	// Iterating through column partitions
	// Column 1 gets mask m1, Column 2 gets mask m2, Column 3 gets remainder
	for _, m1 := range masks3 {
		c1 := getCounts(a, m1)
		s1 := sum(c1)
		
		// Check first X cut
        if s1 > 0 && s1 < n && pointsX[s1-1].x == pointsX[s1].x {
            continue
        }

		for _, m2 := range masks3 {
			if (m1 & m2) != 0 {
				continue
			}
			
			c2 := getCounts(a, m2)
			c3 := getCounts(a, ((1 << 9) - 1) ^ (m1 | m2))

			s2 := sum(c2)
            
            idxX2 := s1 + s2

			// Check second X cut
            if idxX2 > 0 && idxX2 < n && pointsX[idxX2-1].x == pointsX[idxX2].x {
                continue
            }

			// Valid X partitions
            // Calculate xCut1
            var xCut1 float64
            if s1 == 0 {
                if n > 0 { xCut1 = float64(pointsX[0].x) - 0.5 } else { xCut1 = 0.5 }
            } else if s1 == n {
                if n > 0 { xCut1 = float64(pointsX[n-1].x) + 0.5 } else { xCut1 = 0.5 }
            } else {
                xCut1 = float64(pointsX[s1-1].x+pointsX[s1].x) / 2.0
            }

            // Calculate xCut2
            var xCut2 float64
            if idxX2 == 0 {
                if n > 0 { xCut2 = float64(pointsX[0].x) - 0.5 } else { xCut2 = 0.5 }
            } else if idxX2 == n {
                 if n > 0 { xCut2 = float64(pointsX[n-1].x) + 0.5 } else { xCut2 = 0.5 }
            } else {
                xCut2 = float64(pointsX[idxX2-1].x+pointsX[idxX2].x) / 2.0
            }
            
            // Ensure xCut1 < xCut2
            if xCut1 >= xCut2 {
                // This can only happen if s1 == idxX2 (implies s2=0) or if ranges overlap in a way (not possible due to sorting unless same gap)
                // If s1 == idxX2, we are in the same gap.
                // Gap is (pointsX[s1-1].x, pointsX[s1].x). Width >= 1.
                // We picked midpoint.
                // Shift them apart.
                // If s1 == 0 (and idxX2=0), gap is (-inf, pointsX[0].x). We picked pointsX[0].x - 0.5.
                // Make xCut1 smaller.
                if s1 == 0 && idxX2 == 0 {
                    xCut1 -= 1.0
                } else if s1 == n && idxX2 == n {
                    // gap is (pointsX[n-1].x, inf). picked +0.5.
                    // Make xCut2 larger.
                    xCut2 += 1.0
                } else {
                    // In a bounded gap
                    mid := xCut1 // == xCut2
                    xCut1 = mid - 0.1
                    xCut2 = mid + 0.1
                }
            }

			// Partition pointsY into p1, p2, p3
            var thresh1, thresh2 int
            const MAX_INT = 2000000000
            const MIN_INT = -2000000000
            
            if s1 == 0 {
                thresh1 = MIN_INT
            } else if s1 == n {
                thresh1 = MAX_INT
            } else {
                thresh1 = pointsX[s1-1].x
            }

            if idxX2 == 0 {
                thresh2 = MIN_INT
            } else if idxX2 == n {
                thresh2 = MAX_INT
            } else {
                thresh2 = pointsX[idxX2-1].x
            }

			p1 = p1[:0]
			p2 = p2[:0]
			p3 = p3[:0]

			for _, p := range pointsY {
				if p.x <= thresh1 {
					p1 = append(p1, p)
				} else if p.x <= thresh2 {
					p2 = append(p2, p)
				} else {
					p3 = append(p3, p)
				}
			}

			// Get valid Y ranges for each strip
			ranges1 := getValidYRanges(p1, c1)
			ranges2 := getValidYRanges(p2, c2)
			ranges3 := getValidYRanges(p3, c3)

			// Check intersections
			for _, r1 := range ranges1 {
				for _, r2 := range ranges2 {
					iy1_min := max(r1.y1min, r2.y1min)
					iy1_max := min(r1.y1max, r2.y1max)
					iy2_min := max(r1.y2min, r2.y2min)
					iy2_max := min(r1.y2max, r2.y2max)

					if iy1_min < iy1_max && iy2_min < iy2_max {
						for _, r3 := range ranges3 {
							gy1_min := max(iy1_min, r3.y1min)
							gy1_max := min(iy1_max, r3.y1max)
							gy2_min := max(iy2_min, r3.y2min)
							gy2_max := min(iy2_max, r3.y2max)

								if gy1_min < gy1_max && gy2_min < gy2_max {
								// Found solution
								yRes1 := (gy1_min + gy1_max) / 2.0
								yRes2 := (gy2_min + gy2_max) / 2.0
								if yRes1 >= yRes2 {
									yRes1 -= 1e-4
									yRes2 += 1e-4
								}
								fmt.Printf("%.10f %.10f\n", xCut1, xCut2)
								fmt.Printf("%.10f %.10f\n", yRes1, yRes2)
								return
							}
						}
					}
				}
			}
		}
	}

	fmt.Println("-1")
}

func popCount(n int) int {
	c := 0
	for n > 0 {
		if n&1 == 1 {
			c++
		}
		n >>= 1
	}
	return c
}

func getCounts(a []int, mask int) []int {
	res := make([]int, 0, 3)
	for i := 0; i < 9; i++ {
		if (mask & (1 << i)) != 0 {
			res = append(res, a[i])
		}
	}
	return res
}

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

func getValidYRanges(pts []Point, counts []int) []YRange {
	// 3! = 6 permutations
	perms := [][]int{
		{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0},
	}
	res := make([]YRange, 0, 6)
    n := len(pts)
    const INF = 4e9 // Larger than any coordinate range

	for _, p := range perms {
		c1 := counts[p[0]] // bottom count
		c2 := counts[p[1]] // middle count
		
		idx1 := c1
		idx2 := c1 + c2
		
		// Check validity of cuts
        var y1min, y1max float64
        valid1 := true
        
        if idx1 == 0 {
            y1min = -INF
            if n > 0 { y1max = float64(pts[0].y) } else { y1max = INF }
        } else if idx1 == n {
            if n > 0 { y1min = float64(pts[n-1].y) } else { y1min = -INF }
            y1max = INF
        } else {
            if pts[idx1-1].y == pts[idx1].y { valid1 = false }
            y1min = float64(pts[idx1-1].y)
            y1max = float64(pts[idx1].y)
        }

        if !valid1 { continue }
        
        var y2min, y2max float64
        valid2 := true
        
        if idx2 == 0 {
             y2min = -INF
             if n > 0 { y2max = float64(pts[0].y) } else { y2max = INF }
        } else if idx2 == n {
             if n > 0 { y2min = float64(pts[n-1].y) } else { y2min = -INF }
             y2max = INF
        } else {
             if pts[idx2-1].y == pts[idx2].y { valid2 = false }
             y2min = float64(pts[idx2-1].y)
             y2max = float64(pts[idx2].y)
        }
        
		if !valid2 { continue }
        
        // Ensure y1 < y2
        if y1min >= y2max { continue }

		res = append(res, YRange{y1min, y1max, y2min, y2max})
	}
	return res
}

func min(a, b float64) float64 {
	if a < b { return a }
	return b
}

func max(a, b float64) float64 {
	if a > b { return a }
	return b
}