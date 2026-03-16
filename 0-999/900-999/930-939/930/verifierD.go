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

func expectedD(n int, pts [][2]int) int64 {
	if n < 4 {
		return 0
	}
	const INF = int(1e9)
	minX, maxX := INF, -INF
	minY, maxY := INF, -INF
	for _, p := range pts {
		x, y := p[0], p[1]
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	dx := maxX - minX - 1
	dy := maxY - minY - 1
	if dx <= 0 || dy <= 0 {
		return 0
	}
	return int64(dx) * int64(dy)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `6 4 -4 -4 3 2 -3 -4 -5 3 1 5 2
4 1 5 -4 1 -5 5 -1 4
3 -1 -1 3 -4 3 5
6 -5 -3 -5 -4 0 0 2 -2 -2 0 5 4
8 -1 3 -3 2 1 -5 -4 -5 -4 0 -2 -1 3 -4 3 -1
2 -4 3 4 -1
1 -2 -4
3 2 -3 -3 2 -1 -3
8 4 -5 3 1 -5 -4 -4 -1 -4 5 1 -2 -1 1 -2 2
10 3 3 0 -3 1 5 -5 4 1 4 -2 -1 0 3 0 2 5 2 -1 0
5 -3 -3 -2 2 0 5 -5 0 2 4
5 5 -2 -5 -4 2 3 -3 -2 1 4
6 3 -3 -5 2 1 2 2 0 -5 -2 4 -2
2 2 -3 5 -4
9 5 3 1 4 -1 -2 -4 1 -5 3 2 0 -5 -2 4 -1 1 2
9 -1 -4 -5 -4 0 1 5 1 1 5 5 -2 -1 2 1 3 -5 5
7 5 -1 5 2 -3 1 5 1 1 4 4 -1 2 3
6 -1 1 4 5 -3 0 -5 1 2 2 -1 -4
4 0 -2 1 -5 4 0 -3 -3
8 -4 -4 -3 4 -5 -2 5 -3 -5 5 3 5 1 5 1 -3
5 1 2 -1 3 -5 -3 2 -3 -4 5
2 2 -1 -1 -1
7 -5 4 -1 0 2 -4 -5 5 -1 -2 2 5 -1 3
4 1 -4 1 -5 -5 -4 0 -5
5 -1 -5 -5 -5 4 -5 -3 4 -3 2
2 4 -3 -5 3
3 -4 -4 -4 4 5 -4
4 -4 2 1 -4 -4 2 2 3
3 3 -4 -3 0 1 3
3 -4 5 3 -5 3 -2
2 1 0 3 0
9 4 3 1 -1 -2 4 -4 -3 4 1 0 -1 0 3 0 -5 -5 -5
3 -1 5 -2 5 -3 -3
2 2 2 5 -3
5 -1 -3 -5 -3 5 0 0 -2 5 -4
9 2 -3 5 -1 0 -2 5 -4 4 5 -2 0 1 -4 -1 1 3 0
1 5 0
4 1 -5 1 3 -1 -3 -1 -5
5 -4 0 5 1 -5 5 -1 -2 5 -3
5 0 3 5 -4 -3 2 -2 4 4 4
2 1 -1 5 -5
5 -5 -5 4 -3 0 -3 1 -2 4 3
4 -5 5 3 -2 2 3 5 2
10 4 0 0 0 -4 -5 1 0 -3 2 -2 5 5 5 1 4 -5 5 -2 4
8 1 1 -3 -3 5 1 1 2 5 -2 -1 -4 -3 4 -5 -1
3 -4 -3 -1 -3 1 -1
9 5 1 0 3 1 2 2 -4 -1 4 0 -1 -1 -1 1 4 3 -5
8 5 3 3 -2 1 -1 5 1 -5 -5 1 3 -4 -1 -5 -5
8 5 1 3 0 3 2 3 2 3 2 4 -3 4 -5 4 0
4 -1 1 -4 -2 -4 -3 4 2
3 2 0 2 -3 2 -1
8 0 3 -2 -5 0 3 -1 1 3 5 4 -4 -1 -3 3 -3
3 -4 0 -2 -5 -5 2
8 -4 4 3 -5 -3 5 2 4 -3 4 -1 -4 -1 0 -3 2
6 -5 0 -4 -4 4 5 -4 -1 5 -1 4 0
3 -5 5 2 1 1 5
4 -4 1 -1 2 1 5 2 1
3 1 -2 1 5 -2 -1
2 -3 2 0 0
5 2 0 -2 -3 2 -3 0 -2 2 -3
2 3 4 4 2
6 -3 -4 5 0 -3 -1 0 5 -4 -4 1 0
6 -1 0 -5 -2 5 5 -1 5 -3 -3 3 -3
8 -2 -4 3 -1 -4 -5 1 -4 1 2 2 5 -3 -3 -5 -4
3 5 -4 1 4 -4 -1
2 4 4 -4 -2
5 1 -3 3 2 -1 0 5 -3 5 1
5 -3 -3 -4 -2 5 3 0 5 2 5
2 0 3 1 -1
8 -5 -3 0 -3 2 -2 -2 -4 -3 2 1 -4 0 -2 -4 0
3 -2 -2 4 3 0 5
8 0 4 5 2 3 5 -1 -1 -5 4 5 4 -3 -1 -2 3
9 -4 -1 2 -3 -4 -5 2 3 3 3 4 4 3 5 1 1 0 4
6 -2 0 2 -3 3 -5 2 -3 -5 -4 3 -5
3 -4 -5 0 -1 5 -2
2 5 -5 -2 -3
6 1 5 3 4 5 0 1 -2 1 -5 4 -2
10 -5 2 -5 1 -3 -2 -4 4 5 -3 0 -1 -1 2 0 3 1 2 0 -1
9 -1 2 3 0 1 5 0 1 -1 -5 -2 2 -1 1 1 5 -3 -3
8 -2 -3 4 -2 -4 -1 3 3 -3 4 -3 5 2 -4 2 -2
7 -3 -4 -5 3 -5 -4 -4 4 -3 -3 4 -2 5 0
6 -4 -4 1 2 -3 0 4 3 -2 4 -2 4
10 2 4 -1 -4 -3 1 -4 -2 -3 -3 -1 -1 5 0 -3 2 0 3 -4 -4
7 -5 -5 -2 0 1 -2 0 1 4 -4 0 -4 2 4
8 -1 -2 1 -5 0 0 -3 -4 3 -1 2 -5 4 5 -3 -3
9 4 0 5 2 -4 1 -1 2 -1 -5 1 -5 2 -4 -3 -1 0 -1
10 -3 2 2 -5 4 0 -3 1 -5 -1 -3 4 -3 -4 1 3 2 3 5 -3
7 0 -3 3 0 4 -3 -4 -4 -2 -4 -2 2 5 -2
1 1 0
5 0 4 -1 2 3 -5 -1 2 1 -2
8 -5 3 4 1 0 -2 4 0 -1 -3 -3 -1 -2 -4 0 1
3 -1 -1 4 -3 1 -5
2 -3 1 1 -4
2 5 4 -2 3
2 -2 0 2 3
10 -4 -3 -3 4 4 4 -5 2 -1 -2 3 -5 2 -4 -5 5 -5 4 -2 -3
8 5 -2 3 -1 4 4 -4 -3 0 3 -5 4 0 5 -1 3
4 4 -1 -5 4 1 0 -4 0
5 0 -4 -1 4 -3 0 -5 -5 5 4
7 -5 1 -2 1 4 2 -1 -5 -3 4 2 3 -4 5`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+2*n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, 1+2*n, len(parts))
			os.Exit(1)
		}
		pts := make([][2]int, n)
		for i := 0; i < n; i++ {
			x, _ := strconv.Atoi(parts[1+2*i])
			y, _ := strconv.Atoi(parts[2+2*i])
			pts[i] = [2]int{x, y}
		}
		expect := strconv.FormatInt(expectedD(n, pts), 10)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
