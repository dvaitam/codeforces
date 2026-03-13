package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Point struct {
	r, c int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		scanner.Scan()
		v, _ := strconv.Atoi(scanner.Text())
		return v
	}

	if !scanner.Scan() {
		return
	}
	h, _ := strconv.Atoi(scanner.Text())
	w := nextInt()

	grid := make([][]int8, h)
	for i := 0; i < h; i++ {
		grid[i] = make([]int8, w)
		for j := 0; j < w; j++ {
			grid[i][j] = int8(nextInt())
		}
	}

	var erosionKernel []Point
	for dr := -5; dr <= 5; dr++ {
		for dc := -5; dc <= 5; dc++ {
			if dr*dr+dc*dc <= 9 {
				erosionKernel = append(erosionKernel, Point{dr, dc})
			}
		}
	}

	var dilationKernel []Point
	for dr := -6; dr <= 6; dr++ {
		for dc := -6; dc <= 6; dc++ {
			if dr*dr+dc*dc <= 25 {
				dilationKernel = append(dilationKernel, Point{dr, dc})
			}
		}
	}

	visited := make([][]bool, h)
	for i := range visited {
		visited[i] = make([]bool, w)
	}

	var comps [][]Point
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			if grid[r][c] == 1 && !visited[r][c] {
				q := make([]Point, 0, 1024)
				q = append(q, Point{r, c})
				visited[r][c] = true
				head := 0
				for head < len(q) {
					curr := q[head]
					head++
					for dr := -1; dr <= 1; dr++ {
						for dc := -1; dc <= 1; dc++ {
							if dr == 0 && dc == 0 {
								continue
							}
							nr, nc := curr.r+dr, curr.c+dc
							if nr >= 0 && nr < h && nc >= 0 && nc < w {
								if grid[nr][nc] == 1 && !visited[nr][nc] {
									visited[nr][nc] = true
									q = append(q, Point{nr, nc})
								}
							}
						}
					}
				}
				comps = append(comps, q)
			}
		}
	}

	expanded := make([][]int, h)
	for i := range expanded {
		expanded[i] = make([]int, w)
	}

	rayVisited := make([][]bool, h)
	for i := range rayVisited {
		rayVisited[i] = make([]bool, w)
	}

	var sunRays []int
	sun_id := 1

	for _, comp := range comps {
		var core []Point
		for _, p := range comp {
			isCore := true
			for _, d := range erosionKernel {
				nr, nc := p.r+d.r, p.c+d.c
				if nr < 0 || nr >= h || nc < 0 || nc >= w || grid[nr][nc] == 0 {
					isCore = false
					break
				}
			}
			if isCore {
				core = append(core, p)
			}
		}

		if len(core) == 0 {
			continue
		}

		for _, p := range core {
			for _, d := range dilationKernel {
				nr, nc := p.r+d.r, p.c+d.c
				if nr >= 0 && nr < h && nc >= 0 && nc < w {
					expanded[nr][nc] = sun_id
				}
			}
		}

		raysCount := 0
		for _, p := range comp {
			if expanded[p.r][p.c] != sun_id && !rayVisited[p.r][p.c] {
				rayQ := make([]Point, 0, 64)
				rayQ = append(rayQ, p)
				rayVisited[p.r][p.c] = true
				head := 0
				for head < len(rayQ) {
					curr := rayQ[head]
					head++
					for dr := -1; dr <= 1; dr++ {
						for dc := -1; dc <= 1; dc++ {
							if dr == 0 && dc == 0 {
								continue
							}
							nr, nc := curr.r+dr, curr.c+dc
							if nr >= 0 && nr < h && nc >= 0 && nc < w {
								if grid[nr][nc] == 1 && expanded[nr][nc] != sun_id && !rayVisited[nr][nc] {
									rayVisited[nr][nc] = true
									rayQ = append(rayQ, Point{nr, nc})
								}
							}
						}
					}
				}
				if len(rayQ) >= 10 {
					raysCount++
				}
			}
		}
		sunRays = append(sunRays, raysCount)
		sun_id++
	}

	sort.Ints(sunRays)
	fmt.Println(len(sunRays))
	if len(sunRays) > 0 {
		strs := make([]string, len(sunRays))
		for i, v := range sunRays {
			strs[i] = strconv.Itoa(v)
		}
		fmt.Println(strings.Join(strs, " "))
	} else {
		fmt.Println()
	}
}