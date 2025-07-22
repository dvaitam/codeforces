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

func expectedCardio(a []int) string {
	n := len(a)
	total := 0
	for _, x := range a {
		total += x
	}
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if i%2 == 1 {
			y[i] = y[i-1] + a[i-1]
		} else {
			y[i] = y[i-1] - a[i-1]
		}
	}
	minY, maxY := y[0], y[0]
	for i := 1; i <= n; i++ {
		if y[i] < minY {
			minY = y[i]
		}
		if y[i] > maxY {
			maxY = y[i]
		}
	}
	height := maxY - minY
	width := total
	grid := make([][]rune, height)
	for i := range grid {
		row := make([]rune, width)
		for j := range row {
			row[j] = ' '
		}
		grid[i] = row
	}
	currY := 0
	col := 0
	for i := 1; i <= n; i++ {
		step := a[i-1]
		if i%2 == 1 {
			for k := 0; k < step; k++ {
				y0 := currY + k
				row := maxY - 1 - y0
				grid[row][col] = '/'
				col++
			}
			currY += step
		} else {
			for k := 0; k < step; k++ {
				y0 := currY - 1 - k
				row := maxY - 1 - y0
				grid[row][col] = '\\'
				col++
			}
			currY -= step
		}
	}
	var buf bytes.Buffer
	for i := 0; i < height; i++ {
		buf.WriteString(string(grid[i]))
		buf.WriteByte('\n')
	}
	return buf.String()
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
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+n {
			fmt.Fprintf(os.Stderr, "bad test case on line %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(parts[1+i])
		}
		expect := expectedCardio(arr)
		input := fmt.Sprintf("%d\n", n)
		for i, a := range arr {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(a)
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := out.String()
		if got != expect {
			fmt.Printf("test %d failed:\nexpected:\n%s\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
