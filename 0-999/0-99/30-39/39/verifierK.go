package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Embedded correct solver for 39K.
func solveK(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}

	type rect struct {
		r1, r2, c1, c2 int
	}

	var objects []rect
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '*' && !visited[i][j] {
				r2 := i
				for r2+1 < n && grid[r2+1][j] == '*' {
					r2++
				}
				c2 := j
				for c2+1 < m && grid[i][c2+1] == '*' {
					c2++
				}
				for r := i; r <= r2; r++ {
					for c := j; c <= c2; c++ {
						visited[r][c] = true
					}
				}
				objects = append(objects, rect{i + 1, r2 + 1, j + 1, c2 + 1})
			}
		}
	}

	k = len(objects)

	maxInt := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}
	minInt := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	solve := func(S []int) int64 {
		minR, maxR := n+1, 0
		minC, maxC := m+1, 0
		inS := make([]bool, len(objects))

		for _, idx := range S {
			inS[idx] = true
			if objects[idx].r1 < minR {
				minR = objects[idx].r1
			}
			if objects[idx].r2 > maxR {
				maxR = objects[idx].r2
			}
			if objects[idx].c1 < minC {
				minC = objects[idx].c1
			}
			if objects[idx].c2 > maxC {
				maxC = objects[idx].c2
			}
		}

		L_r1, R_r1 := 1, minR
		L_r2, R_r2 := maxR, n
		L_c1, R_c1 := 1, minC
		L_c2, R_c2 := maxC, m

		var topLeft, topRight, bottomLeft, bottomRight []rect

		for i, obj := range objects {
			if inS[i] {
				continue
			}

			rowIntersects := maxInt(obj.r1, minR) <= minInt(obj.r2, maxR)
			colIntersects := maxInt(obj.c1, minC) <= minInt(obj.c2, maxC)

			if rowIntersects && colIntersects {
				return 0
			}

			if obj.r2 < minR && colIntersects {
				if obj.r2+1 > L_r1 {
					L_r1 = obj.r2 + 1
				}
			} else if obj.r1 > maxR && colIntersects {
				if obj.r1-1 < R_r2 {
					R_r2 = obj.r1 - 1
				}
			} else if obj.c2 < minC && rowIntersects {
				if obj.c2+1 > L_c1 {
					L_c1 = obj.c2 + 1
				}
			} else if obj.c1 > maxC && rowIntersects {
				if obj.c1-1 < R_c2 {
					R_c2 = obj.c1 - 1
				}
			} else if obj.r2 < minR && obj.c2 < minC {
				topLeft = append(topLeft, obj)
			} else if obj.r2 < minR && obj.c1 > maxC {
				topRight = append(topRight, obj)
			} else if obj.r1 > maxR && obj.c2 < minC {
				bottomLeft = append(bottomLeft, obj)
			} else if obj.r1 > maxR && obj.c1 > maxC {
				bottomRight = append(bottomRight, obj)
			}
		}

		if L_r1 > R_r1 || L_r2 > R_r2 || L_c1 > R_c1 || L_c2 > R_c2 {
			return 0
		}

		type block struct{ a, b, count int }
		var blocksAB []block
		for r1 := L_r1; r1 <= R_r1; r1++ {
			A := L_c1
			B := R_c2
			for _, obj := range topLeft {
				if obj.r2 >= r1 {
					if obj.c2+1 > A {
						A = obj.c2 + 1
					}
				}
			}
			for _, obj := range topRight {
				if obj.r2 >= r1 {
					if obj.c1-1 < B {
						B = obj.c1 - 1
					}
				}
			}
			if len(blocksAB) > 0 && blocksAB[len(blocksAB)-1].a == A && blocksAB[len(blocksAB)-1].b == B {
				blocksAB[len(blocksAB)-1].count++
			} else {
				blocksAB = append(blocksAB, block{A, B, 1})
			}
		}

		var blocksCD []block
		for r2 := L_r2; r2 <= R_r2; r2++ {
			C := L_c1
			D := R_c2
			for _, obj := range bottomLeft {
				if obj.r1 <= r2 {
					if obj.c2+1 > C {
						C = obj.c2 + 1
					}
				}
			}
			for _, obj := range bottomRight {
				if obj.r1 <= r2 {
					if obj.c1-1 < D {
						D = obj.c1 - 1
					}
				}
			}
			if len(blocksCD) > 0 && blocksCD[len(blocksCD)-1].a == C && blocksCD[len(blocksCD)-1].b == D {
				blocksCD[len(blocksCD)-1].count++
			} else {
				blocksCD = append(blocksCD, block{C, D, 1})
			}
		}

		var ways int64 = 0
		for _, b1 := range blocksAB {
			for _, b2 := range blocksCD {
				w1 := R_c1 - b1.a
				if b2.a > b1.a {
					w1 = R_c1 - b2.a
				}
				w1++
				if w1 <= 0 {
					continue
				}

				w2 := b1.b - L_c2
				if b2.b < b1.b {
					w2 = b2.b - L_c2
				}
				w2++
				if w2 <= 0 {
					continue
				}

				ways += int64(b1.count) * int64(b2.count) * int64(w1) * int64(w2)
			}
		}

		return ways
	}

	var ans int64 = 0

	for i := 0; i < k; i++ {
		ans += solve([]int{i})
	}

	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			ans += solve([]int{i, j})
		}
	}

	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			for l := j + 1; l < k; l++ {
				ans += solve([]int{i, j, l})
			}
		}
	}

	return fmt.Sprintf("%d", ans)
}

func ok(grid [][]byte, r, c int) bool {
	n := len(grid)
	m := len(grid[0])
	for i := r - 1; i <= r+1; i++ {
		for j := c - 1; j <= c+1; j++ {
			if i >= 0 && i < n && j >= 0 && j < m {
				if grid[i][j] == '*' {
					return false
				}
			}
		}
	}
	return true
}

func generateCaseK(rng *rand.Rand) string {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	k := rng.Intn(2) + 1
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	placed := 0
	attempts := 0
	for placed < k && attempts < 100 {
		r := rng.Intn(n)
		c := rng.Intn(m)
		if grid[r][c] == '.' && ok(grid, r, c) {
			grid[r][c] = '*'
			placed++
		}
		attempts++
	}
	k = placed
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseK(rng)
	}
	for i, tc := range cases {
		expect := solveK(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
