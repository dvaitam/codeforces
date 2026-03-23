package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x, y int
}

const testcasesRaw = `100
2
4 2
2 3
1 4
3 4
3
2 3
1 4
2 3
2 3
1 3
4 2
2
2 4
2 3
2 2
4 4
3
4 2
1 4
4 3
3 2
3 4
3 4
2
4 2
2 4
4 1
4 2
3
4 1
3 3
2 2
2 3
2 2
3 2
2
1 1
2 3
3 3
2 1
3
3 1
4 3
3 3
1 1
1 3
4 4
3
2 1
4 4
1 3
3 2
1 1
2 4
3
1 3
3 4
1 1
2 2
3 4
2 1
2
1 3
2 4
1 4
2 3
2
1 2
3 2
2 2
4 4
2
2 1
3 4
2 1
3 3
3
2 4
2 1
1 3
4 1
1 1
3 3
2
3 1
2 3
1 1
4 4
2
3 2
2 2
2 1
1 4
3
2 4
2 2
1 3
2 3
1 4
4 2
3
3 4
2 4
2 2
4 3
4 3
4 2
2
4 1
2 3
3 1
3 4
2
4 2
3 3
1 3
1 4
2
1 1
1 4
3 4
3 1
3
1 4
3 1
3 1
1 3
1 1
1 1
2
4 3
1 1
2 1
2 2
2
4 3
4 2
1 4
1 2
2
4 4
1 2
2 1
2 4
3
4 3
2 1
4 2
4 2
1 3
4 4
3
4 3
2 1
3 2
3 3
2 4
3 3
3
3 3
2 1
3 3
2 1
2 4
2 3
2
4 4
2 3
3 4
1 4
3
1 2
1 3
3 1
4 4
2 1
2 1
3
4 3
4 4
1 3
4 1
1 1
2 2
2
3 4
2 3
4 2
1 3
2
2 2
2 1
1 4
3 3
2
4 4
2 1
2 1
4 3
2
3 1
4 2
3 4
1 4
3
4 2
2 1
1 2
4 3
3 2
1 3
3
1 2
1 4
1 4
3 3
4 2
1 2
2
4 2
3 1
1 2
1 3
3
4 3
3 3
3 2
1 4
1 1
2 2
3
4 2
3 3
3 3
1 1
4 1
4 4
2
1 4
2 4
2 4
1 3
3
3 3
2 1
4 3
2 1
2 4
4 2
2
2 3
2 3
1 3
3 4
2
1 2
3 2
4 3
1 1
3
4 1
2 1
3 4
2 3
4 3
2 1
3
2 1
4 4
3 3
2 2
1 4
2 2
3
2 2
4 2
1 1
2 4
4 4
1 2
3
3 1
2 3
4 4
1 3
1 4
3 2
2
1 3
4 2
3 3
1 1
2
4 2
4 2
4 4
4 2
2
1 4
1 4
3 2
1 3
3
2 4
4 1
2 3
3 3
4 1
4 1
3
3 4
3 4
1 3
2 4
4 3
4 2
3
3 1
3 2
2 3
1 2
1 4
4 3
2
1 4
3 3
4 1
1 1
2
4 1
1 4
4 2
1 3
2
2 4
1 1
1 3
2 4
2
2 2
2 1
1 4
4 3
2
1 1
1 1
2 4
3 1
3
3 1
2 1
3 3
2 2
3 2
2 2
3
1 1
3 4
4 3
2 1
3 4
1 2
2
4 1
3 3
2 1
1 1
2
2 3
1 1
3 3
4 1
2
4 2
3 3
1 2
1 3
3
3 3
1 1
4 3
3 3
1 2
4 3
2
2 1
4 3
1 2
2 2
2
3 4
1 3
2 4
2 4
2
4 2
3 3
3 1
4 1
3
4 4
3 2
2 3
4 4
3 2
2 1
3
2 2
2 1
3 4
1 1
1 3
2 2
3
4 4
3 3
3 3
4 2
4 2
4 2
3
3 3
1 1
2 3
3 2
1 3
3 1
2
4 2
1 4
4 1
2 1
2
2 2
3 3
1 3
3 1
3
1 3
2 4
1 4
2 4
4 3
1 3
2
1 4
4 1
2 1
1 4
2
4 2
4 3
4 1
1 2
2
2 4
2 4
2 1
1 4
3
4 3
4 3
1 2
1 4
4 4
3 4
2
1 1
3 4
4 2
1 2
3
2 2
3 2
4 1
3 4
2 1
3 4
2
4 4
1 3
4 3
3 4
3
4 3
1 2
4 3
2 1
3 1
4 4
2
4 2
4 3
1 1
2 3
2
2 3
1 2
2 2
4 3
3
2 2
4 1
2 4
1 1
3 4
3 2
2
3 3
2 2
3 4
3 3
3
4 4
2 4
2 3
4 4
4 1
3 4
3
1 3
4 3
4 1
1 3
3 4
3 4
2
3 4
4 1
4 3
1 3
2
4 4
2 1
3 2
2 1
3
1 1
3 3
3 1
4 2
1 4
3 4
3
2 3
2 2
2 4
2 2
3 3
2 3
2
1 4
2 2
2 4
2 4
3
1 1
3 2
1 1
3 1
4 3
3 4
3
3 2
3 1
3 4
3 1
2 4
4 3
3
2 2
1 2
1 3
1 2
3 4
3 1
3
1 2
4 1
4 4
2 1
1 3
4 3
2
4 4
2 3
4 3
1 1
3
3 4
2 3
1 1
4 2
4 3
2 2`

func hasRectangle(points []point) bool {
	colRows := make(map[int]map[int]struct{})
	for _, p := range points {
		if colRows[p.y] == nil {
			colRows[p.y] = make(map[int]struct{})
		}
		colRows[p.y][p.x] = struct{}{}
	}
	type pair struct{ a, b int }
	counts := make(map[pair]int)
	for _, rowSet := range colRows {
		rows := make([]int, 0, len(rowSet))
		for r := range rowSet {
			rows = append(rows, r)
		}
		for i := 0; i < len(rows); i++ {
			for j := i + 1; j < len(rows); j++ {
				a, b := rows[i], rows[j]
				if a > b {
					a, b = b, a
				}
				key := pair{a, b}
				counts[key]++
				if counts[key] >= 2 {
					return true
				}
			}
		}
	}
	return false
}

type cell struct{ x, y int }

func solveCase(n int, points []point) (bool, [4]cell) {
	cells := make([]cell, len(points))
	for i, p := range points {
		cells[i] = cell{p.x, p.y}
	}
	rowMap := make(map[int]map[int]int)
	for i, c := range cells {
		if rowMap[c.x] == nil {
			rowMap[c.x] = make(map[int]int)
		}
		rowMap[c.x][c.y] = i
	}
	rows := make([]int, 0, len(rowMap))
	for r := range rowMap {
		rows = append(rows, r)
	}
	sort.Ints(rows)
	for i := 0; i < len(rows); i++ {
		for j := i + 1; j < len(rows); j++ {
			r1, r2 := rows[i], rows[j]
			type colPair struct {
				c, i1, i2 int
			}
			var cols []colPair
			for c1, idx1 := range rowMap[r1] {
				if idx2, ok := rowMap[r2][c1]; ok {
					cols = append(cols, colPair{c1, idx1, idx2})
				}
			}
			if len(cols) < 2 {
				continue
			}
			for a := 0; a < len(cols); a++ {
				for b := a + 1; b < len(cols); b++ {
					i1 := cols[a].i1
					i2 := cols[b].i1
					i3 := cols[a].i2
					i4 := cols[b].i2
					if i1 == i2 || i1 == i3 || i1 == i4 || i2 == i3 || i2 == i4 || i3 == i4 {
						continue
					}
					return true, [4]cell{cells[i1], cells[i2], cells[i3], cells[i4]}
				}
			}
		}
	}
	return false, [4]cell{}
}

func validateYes(points []point, ans [4]cell) bool {
	inSet := make(map[point]struct{})
	for _, p := range points {
		inSet[p] = struct{}{}
	}
	for _, c := range ans {
		if _, ok := inSet[point{c.x, c.y}]; !ok {
			return false
		}
	}
	xs := make(map[int]struct{})
	ys := make(map[int]struct{})
	setCheck := make(map[cell]struct{})
	for _, c := range ans {
		xs[c.x] = struct{}{}
		ys[c.y] = struct{}{}
		setCheck[c] = struct{}{}
	}
	if len(xs) != 2 || len(ys) != 2 || len(setCheck) != 4 {
		return false
	}
	var xv, yv []int
	for x := range xs {
		xv = append(xv, x)
	}
	for y := range ys {
		yv = append(yv, y)
	}
	for _, x := range xv {
		for _, y := range yv {
			if _, ok := inSet[point{x, y}]; !ok {
				return false
			}
		}
	}
	return true
}

func loadCases() ([]int, [][]point, []bool) {
	tokens := strings.Fields(testcasesRaw)
	if len(tokens) == 0 {
		fmt.Println("no embedded testcases")
		os.Exit(1)
	}
	t, err := strconv.Atoi(tokens[0])
	if err != nil {
		fmt.Println("invalid testcase count")
		os.Exit(1)
	}
	var ns []int
	var allPoints [][]point
	var answers []bool
	pos := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos >= len(tokens) {
			fmt.Printf("case %d incomplete\n", caseIdx+1)
			os.Exit(1)
		}
		n, errN := strconv.Atoi(tokens[pos])
		if errN != nil {
			fmt.Printf("case %d invalid n\n", caseIdx+1)
			os.Exit(1)
		}
		pos++
		points := make([]point, 0, 2*n)
		for i := 0; i < 2*n; i++ {
			if pos+1 >= len(tokens) {
				fmt.Printf("case %d missing points\n", caseIdx+1)
				os.Exit(1)
			}
			x, _ := strconv.Atoi(tokens[pos])
			y, _ := strconv.Atoi(tokens[pos+1])
			pos += 2
			points = append(points, point{x, y})
		}
		ns = append(ns, n)
		allPoints = append(allPoints, points)
		answers = append(answers, hasRectangle(points))
	}
	return ns, allPoints, answers
}

func main() {
	// Accept but ignore the binary argument for interface compatibility
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	ns, allPoints, exists := loadCases()
	for idx := range allPoints {
		points := allPoints[idx]
		hasRect := exists[idx]
		n := ns[idx]

		found, ans := solveCase(n, points)

		if hasRect {
			if !found {
				fmt.Printf("case %d: embedded solver failed to find rectangle that exists\n", idx+1)
				os.Exit(1)
			}
			if !validateYes(points, ans) {
				fmt.Printf("case %d: embedded solver produced invalid rectangle\n", idx+1)
				os.Exit(1)
			}
		} else {
			if found {
				fmt.Printf("case %d: embedded solver found rectangle but none should exist\n", idx+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(allPoints))
}
