package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type state struct {
	ready    int64
	prevLane int8
	shot     bool
	shotTime int64
	enterPos int
}

const inf int64 = 1<<62 - 1

func minState(a, b state) state {
	if a.ready <= b.ready {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var m1, m2 int
	var t int
	if _, err := fmt.Fscan(in, &n, &m1, &m2, &t); err != nil {
		return
	}

	eventsMap := make(map[int]uint8, m1+m2)

	for i := 0; i < m1; i++ {
		var x int
		fmt.Fscan(in, &x)
		eventsMap[x] |= 1
	}
	for i := 0; i < m2; i++ {
		var x int
		fmt.Fscan(in, &x)
		eventsMap[x] |= 2
	}

	xs := make([]int, 0, len(eventsMap)+2)
	xs = append(xs, 0)
	for x := range eventsMap {
		xs = append(xs, x)
	}
	xs = append(xs, n+1)
	sort.Ints(xs)

	type event struct {
		x    int
		mask uint8
	}

	events := make([]event, 0, len(xs))
	var last int = -1
	for _, x := range xs {
		if x == last {
			continue
		}
		last = x
		events = append(events, event{x: x, mask: eventsMap[x]})
	}

	m := len(events)
	dp := make([][2]state, m)

	initialReady := int64(t)
	for lane := 0; lane < 2; lane++ {
		dp[0][lane] = state{
			ready:    initialReady,
			prevLane: -1,
			enterPos: 0,
		}
	}

	for i := 1; i < m; i++ {
		xPrev := events[i-1].x
		xCur := events[i].x
		gap := xCur - xPrev
		for lane := 0; lane < 2; lane++ {
			best := state{ready: inf}
			for prev := 0; prev < 2; prev++ {
				prevState := dp[i-1][prev]
				if prevState.ready >= inf {
					continue
				}
				if gap == 1 && prev != lane {
					continue
				}
				earliest := xPrev
				if prev != lane {
					if gap == 1 {
						continue
					}
					if events[i-1].mask&(1<<lane) != 0 {
						earliest = xPrev + 1
						if earliest >= xCur {
							continue
						}
					}
				}
				cur := state{
					prevLane: int8(prev),
					enterPos: earliest,
				}
				readyBefore := prevState.ready
				if events[i].mask&(1<<lane) != 0 {
					shotTime := readyBefore
					if shotTime < int64(earliest) {
						shotTime = int64(earliest)
					}
					if shotTime > int64(xCur-1) {
						continue
					}
					cur.ready = shotTime + int64(t)
					cur.shot = true
					cur.shotTime = shotTime
				} else {
					cur.ready = readyBefore
				}
				if cur.ready < best.ready {
					best = cur
				}
			}
			dp[i][lane] = best
		}
	}

	lastIdx := m - 1
	finalLane := -1
	bestReady := inf
	for lane := 0; lane < 2; lane++ {
		if dp[lastIdx][lane].ready < bestReady {
			bestReady = dp[lastIdx][lane].ready
			finalLane = lane
		}
	}
	if finalLane == -1 || bestReady >= inf {
		fmt.Println("No")
		return
	}

	laneAt := make([]int, m)
	lane := finalLane
	for i := lastIdx; i >= 0; i-- {
		laneAt[i] = lane
		prev := dp[i][lane].prevLane
		lane = int(prev)
	}

	transitions := make([]int, 0)
	for i := 1; i < m; i++ {
		if laneAt[i] != laneAt[i-1] {
			transitions = append(transitions, dp[i][laneAt[i]].enterPos)
		}
	}

	type shotInfo struct {
		x int64
		y int
	}

	shots := make([]shotInfo, 0)
	for i := 1; i < m; i++ {
		st := dp[i][laneAt[i]]
		if st.ready >= inf {
			continue
		}
		if st.shot {
			shots = append(shots, shotInfo{x: st.shotTime, y: laneAt[i] + 1})
		}
	}

	fmt.Println("Yes")
	fmt.Println(len(transitions))
	if len(transitions) > 0 {
		for i, v := range transitions {
			if i > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%d", v)
		}
		fmt.Println()
	} else {
		fmt.Println()
	}
	fmt.Println(len(shots))
	for _, sh := range shots {
		fmt.Printf("%d %d\n", sh.x, sh.y)
	}
}
