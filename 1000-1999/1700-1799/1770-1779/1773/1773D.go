package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const (
	modulus       int64 = 1000000007
	limitAnswer         = 1000000
	minExactCells       = 2002
)

var dirs = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

type cell struct {
	x int
	y int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(in, &row)
		board[i] = []byte(row)
	}

	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}

	var components [][]cell
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if board[i][j] != '.' || visited[i][j] {
				continue
			}
			stack := []cell{{i, j}}
			visited[i][j] = true
			var comp []cell
			for len(stack) > 0 {
				last := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				comp = append(comp, last)
				for _, d := range dirs {
					ni := last.x + d[0]
					nj := last.y + d[1]
					if ni >= 0 && ni < n && nj >= 0 && nj < m && board[ni][nj] == '.' && !visited[ni][nj] {
						visited[ni][nj] = true
						stack = append(stack, cell{ni, nj})
					}
				}
			}
			components = append(components, comp)
		}
	}

	totalCells := 0
	for _, comp := range components {
		totalCells += len(comp)
	}

	if totalCells >= minExactCells {
		fmt.Println(limitAnswer)
		return
	}

	totalPairs := int64(totalCells) * int64(totalCells-1) / 2
	var goodPairs int64
	rng := rand.New(rand.NewSource(1))

	for _, comp := range components {
		if len(comp) < 2 {
			continue
		}
		goodPairs += countGoodPairs(comp, n, m, rng)
	}

	bad := totalPairs - goodPairs
	if bad > limitAnswer {
		bad = limitAnswer
	}
	if bad < 0 {
		bad = 0
	}
	fmt.Println(bad)
}

func countGoodPairs(comp []cell, height, width int, rng *rand.Rand) int64 {
	whiteIndex := make(map[int]int)
	blackIndex := make(map[int]int)
	whites := make([]cell, 0, (len(comp)+1)/2)
	blacks := make([]cell, 0, (len(comp)+1)/2)

	for _, c := range comp {
		id := c.x*width + c.y
		if (c.x+c.y)&1 == 0 {
			whiteIndex[id] = len(whites)
			whites = append(whites, c)
		} else {
			blackIndex[id] = len(blacks)
			blacks = append(blacks, c)
		}
	}

	if len(whites) == 0 || len(whites) != len(blacks) {
		return 0
	}

	size := len(whites)
	if size == 0 {
		return 0
	}

	type edge struct {
		w int
		b int
	}
	edges := make([]edge, 0, size*4)
	for wi, c := range whites {
		for _, d := range dirs {
			nx := c.x + d[0]
			ny := c.y + d[1]
			if nx < 0 || nx >= height || ny < 0 || ny >= width {
				continue
			}
			id := nx*width + ny
			if bi, ok := blackIndex[id]; ok {
				edges = append(edges, edge{wi, bi})
			}
		}
	}

	matrix := make([][]int64, size)
	buf := make([]int64, size*size)
	for i := 0; i < size; i++ {
		matrix[i] = buf[i*size : (i+1)*size]
	}

	const maxAttempts = 60
	for attempt := 0; attempt < maxAttempts; attempt++ {
		for _, e := range edges {
			matrix[e.w][e.b] = rng.Int63n(modulus-1) + 1
		}
		inv, ok := invertMatrix(matrix)
		if ok {
			var count int64
			for _, row := range inv {
				for _, val := range row {
					if val != 0 {
						count++
					}
				}
			}
			return count
		}
	}
	// extremely unlikely fallback: assume worst case (no good pairs)
	return 0
}

func invertMatrix(mat [][]int64) ([][]int64, bool) {
	n := len(mat)
	aug := make([][]int64, n)
	for i := 0; i < n; i++ {
		row := make([]int64, 2*n)
		copy(row, mat[i])
		row[n+i] = 1
		aug[i] = row
	}

	for col := 0; col < n; col++ {
		pivot := col
		for pivot < n && aug[pivot][col] == 0 {
			pivot++
		}
		if pivot == n {
			return nil, false
		}
		if pivot != col {
			aug[pivot], aug[col] = aug[col], aug[pivot]
		}
		invPivot := modInverse(aug[col][col])
		for j := 0; j < 2*n; j++ {
			aug[col][j] = aug[col][j] * invPivot % modulus
		}
		for row := 0; row < n; row++ {
			if row == col {
				continue
			}
			factor := aug[row][col]
			if factor == 0 {
				continue
			}
			for j := 0; j < 2*n; j++ {
				val := aug[row][j] - factor*aug[col][j]%modulus
				val %= modulus
				if val < 0 {
					val += modulus
				}
				aug[row][j] = val
			}
		}
	}

	inv := make([][]int64, n)
	for i := 0; i < n; i++ {
		row := make([]int64, n)
		copy(row, aug[i][n:])
		inv[i] = row
	}
	return inv, true
}

func modInverse(x int64) int64 {
	return modPow(x, modulus-2)
}

func modPow(a, e int64) int64 {
	res := int64(1)
	base := a % modulus
	for e > 0 {
		if e&1 == 1 {
			res = res * base % modulus
		}
		base = base * base % modulus
		e >>= 1
	}
	return res
}
