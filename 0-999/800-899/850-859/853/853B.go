package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type flight struct {
	day  int
	city int
	cost int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	arrivals := make([]flight, 0, m)
	departures := make([]flight, 0, m)
	maxDay := 0
	for i := 0; i < m; i++ {
		var d, f, t int
		var c int64
		fmt.Fscan(in, &d, &f, &t, &c)
		if d > maxDay {
			maxDay = d
		}
		if t == 0 {
			arrivals = append(arrivals, flight{d, f, c})
		} else if f == 0 {
			departures = append(departures, flight{d, t, c})
		}
	}
	maxIndex := maxDay + k + 5
	const INF int64 = 1 << 60
	arrCost := make([]int64, maxIndex+2)
	depCost := make([]int64, maxIndex+2)
	for i := range arrCost {
		arrCost[i] = INF
		depCost[i] = INF
	}

	sort.Slice(arrivals, func(i, j int) bool { return arrivals[i].day < arrivals[j].day })
	bestArr := make([]int64, n+1)
	for i := range bestArr {
		bestArr[i] = INF
	}
	sum := int64(0)
	cntMissing := n
	idx := 0
	for day := 0; day <= maxIndex; day++ {
		for idx < len(arrivals) && arrivals[idx].day == day {
			f := arrivals[idx]
			if bestArr[f.city] == INF {
				bestArr[f.city] = f.cost
				sum += f.cost
				cntMissing--
			} else if f.cost < bestArr[f.city] {
				sum += f.cost - bestArr[f.city]
				bestArr[f.city] = f.cost
			}
			idx++
		}
		if cntMissing == 0 {
			arrCost[day] = sum
		}
	}
	// propagate cost forward for days without updates
	for day := 1; day <= maxIndex; day++ {
		if arrCost[day] > arrCost[day-1] {
			arrCost[day] = arrCost[day-1]
		}
	}

	sort.Slice(departures, func(i, j int) bool { return departures[i].day > departures[j].day })
	bestDep := make([]int64, n+1)
	for i := range bestDep {
		bestDep[i] = INF
	}
	sum = 0
	cntMissing = n
	idx = 0
	for day := maxIndex; day >= 0; day-- {
		for idx < len(departures) && departures[idx].day == day {
			f := departures[idx]
			if bestDep[f.city] == INF {
				bestDep[f.city] = f.cost
				sum += f.cost
				cntMissing--
			} else if f.cost < bestDep[f.city] {
				sum += f.cost - bestDep[f.city]
				bestDep[f.city] = f.cost
			}
			idx++
		}
		if cntMissing == 0 {
			depCost[day] = sum
		}
	}
	for day := maxIndex - 1; day >= 0; day-- {
		if depCost[day] > depCost[day+1] {
			depCost[day] = depCost[day+1]
		}
	}

	ans := INF
	for day := 0; day <= maxDay; day++ {
		dIdx := day + k + 1
		if dIdx > maxIndex {
			break
		}
		if arrCost[day] == INF || depCost[dIdx] == INF {
			continue
		}
		total := arrCost[day] + depCost[dIdx]
		if total < ans {
			ans = total
		}
	}
	if ans == INF {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
