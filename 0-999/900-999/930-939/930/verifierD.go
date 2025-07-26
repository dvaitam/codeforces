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
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
