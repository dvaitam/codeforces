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

type DSU struct {
	parent []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{parent: p}
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	fx, fy := d.Find(x), d.Find(y)
	if fx != fy {
		d.parent[fx] = fy
	}
}

func solveC(n int, p []int) int {
	d := NewDSU(n)
	for i := 1; i <= n; i++ {
		d.Union(i, p[i-1])
	}
	seen := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		seen[d.Find(i)] = struct{}{}
	}
	return len(seen)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Println("could not open testcasesC.txt:", err)
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
		fields := strings.Fields(line)
		n := atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Printf("bad test %d\n", idx)
			os.Exit(1)
		}
		p := make([]int, n)
		for i := 0; i < n; i++ {
			p[i] = atoi(fields[1+i])
		}
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		expected := fmt.Sprintf("%d", solveC(n, p))

		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func atoi(s string) int { v, _ := strconv.Atoi(s); return v }
