package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateTest() string {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	k := rand.Intn(100) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			val := rand.Intn(20) + 1
			if j > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", val)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func check(input, output string) error {
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	firstLine := strings.TrimSpace(lines[0])
	if firstLine == "NO" {
		return nil 
	}
	if firstLine != "YES" {
		return fmt.Errorf("expected YES or NO, got %q", firstLine)
	}

	// Parse input to verify validity
	inLines := strings.Fields(input)
	n, _ := strconv.Atoi(inLines[0])
	m, _ := strconv.Atoi(inLines[1])
	k, _ := strconv.ParseInt(inLines[2], 10, 64)
	
	grid := make([][]int, n)
	idx := 3
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j], _ = strconv.Atoi(inLines[idx])
			idx++
		}
	}

	// Parse output grid
	if len(lines) < 1+n {
		return fmt.Errorf("expected %d lines of grid, got %d", n, len(lines)-1)
	}
	
	resGrid := make([][]int64, n)
	var sum int64
	var targetVal int64 = -1
	var startR, startC int = -1, -1
	count := 0

	for i := 0; i < n; i++ {
		resGrid[i] = make([]int64, m)
		parts := strings.Fields(lines[1+i])
		if len(parts) != m {
			return fmt.Errorf("row %d has %d cols, expected %d", i, len(parts), m)
		}
		for j := 0; j < m; j++ {
			val, _ := strconv.ParseInt(parts[j], 10, 64)
			resGrid[i][j] = val
			sum += val
			if val > 0 {
				count++
				if targetVal == -1 {
					targetVal = val
				} else if targetVal != val {
					return fmt.Errorf("multiple non-zero values found: %d and %d", targetVal, val)
				}
				if int64(grid[i][j]) < val {
					return fmt.Errorf("cell (%d,%d) has original value %d, but result uses %d", i, j, grid[i][j], val)
				}
				if startR == -1 {
					startR, startC = i, j
				}
			}
		}
	}

	if sum != k {
		return fmt.Errorf("sum of hay is %d, expected %d", sum, k)
	}
	
	// Check connectivity
	if count > 0 {
		visited := make([][]bool, n)
		for i := 0; i < n; i++ {
			visited[i] = make([]bool, m)
		}
		
		q := []int{startR*m + startC}
		visited[startR][startC] = true
		seenCount := 0
		
		dr := []int{0, 0, 1, -1}
		dc := []int{1, -1, 0, 0}

		for len(q) > 0 {
			curr := q[0]
			q = q[1:]
			seenCount++
			r, c := curr/m, curr%m
			
			for i := 0; i < 4; i++ {
				nr, nc := r+dr[i], c+dc[i]
				if nr >= 0 && nr < n && nc >= 0 && nc < m {
					if !visited[nr][nc] && resGrid[nr][nc] > 0 {
						visited[nr][nc] = true
						q = append(q, nr*m + nc)
					}
				}
			}
		}
		
		if seenCount != count {
			return fmt.Errorf("connected component size is %d, but total non-zero cells are %d", seenCount, count)
		}
	}
	
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 100; i++ {
		input := generateTest()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(input, output); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
