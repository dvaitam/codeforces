package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Station struct {
	t int
	x int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var e, s int
	var n, m int
	if _, err := fmt.Fscan(in, &e, &s, &n, &m); err != nil {
		return
	}
	stations := make([]Station, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &stations[i].t, &stations[i].x)
	}
	sort.Slice(stations, func(i, j int) bool { return stations[i].x < stations[j].x })
	starts := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &starts[i])
	}

	// For each start position, run a simple greedy simulation.
	for _, f := range starts {
		pos := f
		fuel := s
		fuelType := 3
		c1, c2 := 0, 0
		idx := sort.Search(len(stations), func(i int) bool { return stations[i].x > pos })
		unreachable := false
		for pos < e {
			distToE := e - pos
			if distToE <= fuel {
				// go directly
				if fuelType == 1 {
					c1 += distToE
				} else if fuelType == 2 {
					c2 += distToE
				}
				pos = e
				break
			}
			// find reachable station ahead
			best := -1
			bestType := 0
			bestDist := 0
			j := idx
			for j < len(stations) && stations[j].x-pos <= fuel {
				d := stations[j].x - pos
				t := stations[j].t
				if best == -1 || t > bestType || (t == bestType && d < bestDist) {
					best = j
					bestType = t
					bestDist = d
					if bestType == 3 && bestDist == d {
						break
					}
				}
				j++
			}
			if best == -1 {
				unreachable = true
				break
			}
			// travel to station
			dist := stations[best].x - pos
			if fuelType == 1 {
				c1 += dist
			} else if fuelType == 2 {
				c2 += dist
			}
			fuel -= dist
			pos = stations[best].x
			fuel = s
			fuelType = stations[best].t
			idx = best + 1
		}
		if unreachable {
			fmt.Fprintln(out, "-1 -1")
		} else {
			fmt.Fprintf(out, "%d %d\n", c1, c2)
		}
	}
}
