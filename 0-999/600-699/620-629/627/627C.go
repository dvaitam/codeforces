package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Station struct {
	x int64
	p int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var d, n int64
	var m int
	if _, err := fmt.Fscan(reader, &d, &n, &m); err != nil {
		return
	}
	stations := make([]Station, m+2)
	// first and last stations
	stations[0] = Station{0, 0}
	for i := 0; i < m; i++ {
		var x, p int64
		fmt.Fscan(reader, &x, &p)
		stations[i+1] = Station{x, p}
	}
	stations[m+1] = Station{d, 0}
	sort.Slice(stations, func(i, j int) bool { return stations[i].x < stations[j].x })

	// feasibility check
	for i := 0; i < len(stations)-1; i++ {
		if stations[i+1].x-stations[i].x > n {
			fmt.Println(-1)
			return
		}
	}

	// next cheaper station index
	next := make([]int, len(stations))
	stack := []int{}
	for i := len(stations) - 1; i >= 0; i-- {
		for len(stack) > 0 && stations[stack[len(stack)-1]].p >= stations[i].p {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = -1
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	fuel := n
	var cost int64
	for i := 0; i < len(stations)-1; i++ {
		dist := stations[i+1].x - stations[i].x
		target := n
		if nxt := next[i]; nxt != -1 {
			if stations[nxt].x-stations[i].x <= n {
				target = stations[nxt].x - stations[i].x
			}
		}
		if target < dist {
			target = dist
		}
		if fuel < target {
			add := target - fuel
			cost += add * stations[i].p
			fuel += add
		}
		fuel -= dist
	}
	fmt.Println(cost)
}
