package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Rook struct {
	x, y  int
	color int
}

func parseOutput(n int, out string) ([]Rook, [][]int, error) {
	r := bufio.NewReader(strings.NewReader(out))
	var rooks []Rook
	colorIdx := make([][]int, n+1)
	coords := make(map[[2]int]bool)
	total := 0
	for c := 1; c <= n; c++ {
		var cnt int
		if _, err := fmt.Fscan(r, &cnt); err != nil {
			return nil, nil, fmt.Errorf("failed to read count for color %d: %v", c, err)
		}
		if cnt <= 0 || cnt > 5000 {
			return nil, nil, fmt.Errorf("invalid count for color %d", c)
		}
		for i := 0; i < cnt; i++ {
			var x, y int
			if _, err := fmt.Fscan(r, &x, &y); err != nil {
				return nil, nil, fmt.Errorf("failed to read rook for color %d: %v", c, err)
			}
			if x < 1 || x > 1_000_000_000 || y < 1 || y > 1_000_000_000 {
				return nil, nil, fmt.Errorf("coordinates out of range: %d %d", x, y)
			}
			key := [2]int{x, y}
			if coords[key] {
				return nil, nil, fmt.Errorf("duplicate cell %d %d", x, y)
			}
			coords[key] = true
			rooks = append(rooks, Rook{x: x, y: y, color: c})
			colorIdx[c] = append(colorIdx[c], len(rooks)-1)
		}
		total += cnt
		if total > 5000 {
			return nil, nil, fmt.Errorf("total rooks exceed limit")
		}
	}
	if _, err := fmt.Fscan(r, new(int)); err == nil {
		return nil, nil, fmt.Errorf("extra output")
	}
	return rooks, colorIdx, nil
}

func isConnected(indices []int, rooks []Rook, rowMap, colMap map[int][]int) bool {
	if len(indices) == 0 {
		return false
	}
	allowed := make(map[int]bool, len(indices))
	for _, idx := range indices {
		allowed[idx] = true
	}
	visited := make(map[int]bool)
	q := []int{indices[0]}
	visited[indices[0]] = true
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		row := rooks[u].x
		for _, v := range rowMap[row] {
			if allowed[v] && !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
		col := rooks[u].y
		for _, v := range colMap[col] {
			if allowed[v] && !visited[v] {
				visited[v] = true
				q = append(q, v)
			}
		}
	}
	return len(visited) == len(indices)
}

func verifyCase(n int, edges [][2]int, output string) error {
	rooks, colorIdx, err := parseOutput(n, output)
	if err != nil {
		return err
	}
	rowMap := make(map[int][]int)
	colMap := make(map[int][]int)
	for i, r := range rooks {
		rowMap[r.x] = append(rowMap[r.x], i)
		colMap[r.y] = append(colMap[r.y], i)
	}
	for c := 1; c <= n; c++ {
		if !isConnected(colorIdx[c], rooks, rowMap, colMap) {
			return fmt.Errorf("color %d rooks not connected", c)
		}
	}
	adj := make([][]bool, n+1)
	for i := range adj {
		adj[i] = make([]bool, n+1)
	}
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a][b] = true
		adj[b][a] = true
	}
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			inds := append(append([]int{}, colorIdx[a]...), colorIdx[b]...)
			conn := isConnected(inds, rooks, rowMap, colMap)
			if adj[a][b] && !conn {
				return fmt.Errorf("colors %d and %d should be connected", a, b)
			}
			if !adj[a][b] && conn {
				return fmt.Errorf("colors %d and %d should not be connected", a, b)
			}
		}
	}
	return nil
}

func runSolution(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		idx++
		parts := strings.Fields(line1)
		if len(parts) != 2 {
			fmt.Printf("invalid header on test %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		mVal, _ := strconv.Atoi(parts[1])
		if !scanner.Scan() {
			fmt.Printf("missing edge line for test %d\n", idx)
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scanner.Text())
		nums := []string{}
		if line2 != "" {
			nums = strings.Fields(line2)
		}
		if len(nums) != 2*mVal {
			fmt.Printf("wrong number of edge values in test %d\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, mVal)
		for j := 0; j < mVal; j++ {
			x, _ := strconv.Atoi(nums[2*j])
			y, _ := strconv.Atoi(nums[2*j+1])
			edges[j] = [2]int{x, y}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", nVal, mVal))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runSolution(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := verifyCase(nVal, edges, out); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
