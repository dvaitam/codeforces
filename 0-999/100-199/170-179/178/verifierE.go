package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	grid := make([]byte, n*n)
	for i := 0; i < n*n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if v != 0 {
			grid[i] = 1
		}
	}
	visited := make([]bool, n*n)
	dirs := []int{-n, n, -1, 1}
	circles, squares := 0, 0
	for idx := 0; idx < n*n; idx++ {
		if grid[idx] == 1 && !visited[idx] {
			q := []int{idx}
			visited[idx] = true
			sumX, sumY, sumXXYY := 0.0, 0.0, 0.0
			area := 0
			for qi := 0; qi < len(q); qi++ {
				u := q[qi]
				x := u / n
				y := u % n
				area++
				fx := float64(x)
				fy := float64(y)
				sumX += fx
				sumY += fy
				sumXXYY += fx*fx + fy*fy
				for _, d := range dirs {
					v := u + d
					if d == -1 && y == 0 {
						continue
					}
					if d == 1 && y == n-1 {
						continue
					}
					if v < 0 || v >= len(grid) {
						continue
					}
					if grid[v] == 1 && !visited[v] {
						visited[v] = true
						q = append(q, v)
					}
				}
			}
			if area < 100 {
				continue
			}
			areaF := float64(area)
			xc := sumX / areaF
			yc := sumY / areaF
			E := sumXXYY/areaF - (xc*xc + yc*yc)
			Nc := 2 * math.Pi * E / areaF
			if math.Abs(Nc-1) < math.Abs(Nc-math.Pi/3) {
				circles++
			} else {
				squares++
			}
		}
	}
	return fmt.Sprintf("%d %d", circles, squares)
}

type testE struct{ input, expect string }

func genGrid() [][]int {
	n := 30
	g := make([][]int, n)
	for i := range g {
		g[i] = make([]int, n)
	}
	for i := 2; i < 12; i++ {
		for j := 2; j < 12; j++ {
			g[i][j] = 1
		}
	}
	cx, cy, r := 20, 15, 6
	for i := cx - r; i <= cx+r; i++ {
		for j := cy - r; j <= cy+r; j++ {
			if i >= 0 && j >= 0 && i < n && j < n {
				dx := i - cx
				dy := j - cy
				if dx*dx+dy*dy <= r*r {
					g[i][j] = 1
				}
			}
		}
	}
	return g
}

func gridInput(g [][]int) string {
	var sb strings.Builder
	n := len(g)
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, row := range g {
		for j, v := range row {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genTests() []testE {
	g := genGrid()
	input := gridInput(g)
	expect := solveE(input)
	tests := make([]testE, 100)
	for i := 0; i < 100; i++ {
		tests[i] = testE{input, expect}
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\nexpect:%s\nactual:%s\n", i+1, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
