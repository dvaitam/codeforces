package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func solveB(n, m int, grid [][]byte) int {
	// BFS on bipartite graph of rows and columns connected by '#' cells.
	// Answer is the BFS distance from row 0 to row n-1 in this bipartite graph.
	rowHasCol := make([][]int, n)
	colHasRow := make([][]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				rowHasCol[i] = append(rowHasCol[i], j)
				colHasRow[j] = append(colHasRow[j], i)
			}
		}
	}

	type State struct {
		isRow bool
		idx   int
	}
	visRow := make([]bool, n)
	visCol := make([]bool, m)
	distRow := make([]int, n)
	distCol := make([]int, m)

	queue := []State{{isRow: true, idx: 0}}
	visRow[0] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr.isRow {
			if curr.idx == n-1 {
				return distRow[n-1]
			}
			for _, c := range rowHasCol[curr.idx] {
				if !visCol[c] {
					visCol[c] = true
					distCol[c] = distRow[curr.idx] + 1
					queue = append(queue, State{isRow: false, idx: c})
				}
			}
		} else {
			for _, r := range colHasRow[curr.idx] {
				if !visRow[r] {
					visRow[r] = true
					distRow[r] = distCol[curr.idx] + 1
					queue = append(queue, State{isRow: true, idx: r})
				}
			}
		}
	}
	return -1
}

func runCase(bin string, n, m int, grid [][]byte) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		input.Write(grid[i])
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(&out, &got); err != nil {
		return fmt.Errorf("parse error: %v", err)
	}
	want := solveB(n, m, grid)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(10) + 2
		m := rng.Intn(10) + 2
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(3) == 0 {
					row[j] = '#'
				} else {
					row[j] = '.'
				}
			}
			grid[i] = row
		}
		if err := runCase(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
