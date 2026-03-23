package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	rowObs := make(map[int][]int)
	colObs := make(map[int][]int)
	obsSet := make(map[[2]int]bool)

	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		rowObs[x] = append(rowObs[x], y)
		colObs[y] = append(colObs[y], x)
		obsSet[[2]int{x, y}] = true
	}

	for _, v := range rowObs {
		sort.Ints(v)
	}
	for _, v := range colObs {
		sort.Ints(v)
	}

	total := n*m - k
	x, y := 1, 1
	if obsSet[[2]int{1, 1}] {
		fmt.Println("No")
		return
	}

	visited := 1
	if visited == total {
		fmt.Println("Yes")
		return
	}

	minX, maxX := 1, n
	minY, maxY := 1, m
	dir := 0 // 0: right, 1: down, 2: left, 3: up
	turns := 0

	for {
		var steps int
		var newX, newY int

		switch dir {
		case 0: // right
			limit := maxY
			if obs, ok := rowObs[x]; ok {
				idx := sort.Search(len(obs), func(i int) bool { return obs[i] > y })
				if idx < len(obs) && obs[idx]-1 < limit {
					limit = obs[idx] - 1
				}
			}
			if limit > maxY {
				limit = maxY
			}
			steps = limit - y
			newX, newY = x, limit
		case 1: // down
			limit := maxX
			if obs, ok := colObs[y]; ok {
				idx := sort.Search(len(obs), func(i int) bool { return obs[i] > x })
				if idx < len(obs) && obs[idx]-1 < limit {
					limit = obs[idx] - 1
				}
			}
			if limit > maxX {
				limit = maxX
			}
			steps = limit - x
			newX, newY = limit, y
		case 2: // left
			limit := minY
			if obs, ok := rowObs[x]; ok {
				idx := sort.Search(len(obs), func(i int) bool { return obs[i] >= y })
				if idx > 0 && obs[idx-1]+1 > limit {
					limit = obs[idx-1] + 1
				}
			}
			if limit < minY {
				limit = minY
			}
			steps = y - limit
			newX, newY = x, limit
		case 3: // up
			limit := minX
			if obs, ok := colObs[y]; ok {
				idx := sort.Search(len(obs), func(i int) bool { return obs[i] >= x })
				if idx > 0 && obs[idx-1]+1 > limit {
					limit = obs[idx-1] + 1
				}
			}
			if limit < minX {
				limit = minX
			}
			steps = x - limit
			newX, newY = limit, y
		}

		if steps < 0 {
			steps = 0
		}

		if steps == 0 {
			dir = (dir + 1) % 4
			turns++
			if turns == 4 {
				break
			}
			continue
		}

		turns = 0
		x, y = newX, newY
		visited += steps

		switch dir {
		case 0:
			minX++
		case 1:
			maxY--
		case 2:
			maxX--
		case 3:
			minY++
		}

		dir = (dir + 1) % 4

		if visited == total {
			fmt.Println("Yes")
			return
		}

		if visited > total {
			break
		}
	}

	fmt.Println("No")
}
