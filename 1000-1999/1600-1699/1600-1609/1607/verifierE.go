package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func solve(n, m int, s string) (int, int) {
	minX, maxX := 0, 0
	minY, maxY := 0, 0
	x, y := 0, 0
	startRow, startCol := 1, 1
	for _, c := range s {
		nx, ny := x, y
		switch c {
		case 'L':
			ny--
		case 'R':
			ny++
		case 'U':
			nx--
		case 'D':
			nx++
		}
		newMinX := min(minX, nx)
		newMaxX := max(maxX, nx)
		newMinY := min(minY, ny)
		newMaxY := max(maxY, ny)
		if newMaxX-newMinX+1 > n || newMaxY-newMinY+1 > m {
			break
		}
		x, y = nx, ny
		minX, maxX = newMinX, newMaxX
		minY, maxY = newMinY, newMaxY
		startRow = 1 - minX
		startCol = 1 - minY
	}
	return startRow, startCol
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func randCase() (int, int, string) {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	length := rand.Intn(20) + 1
	b := make([]byte, length)
	dirs := "LRUD"
	for i := range b {
		b[i] = dirs[rand.Intn(4)]
	}
	return n, m, string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	rand.Seed(5)
	ns := []int{1, 2, 3}
	ms := []int{1, 2, 3}
	ss := []string{"L", "RRUU", "UDLR"}
	cases := make([]struct {
		n, m int
		s    string
	}, 0, 100)
	for i := 0; i < len(ns); i++ {
		cases = append(cases, struct {
			n, m int
			s    string
		}{ns[i], ms[i], ss[i]})
	}
	for len(cases) < 100 {
		n, m, s := randCase()
		cases = append(cases, struct {
			n, m int
			s    string
		}{n, m, s})
	}

	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d\n%s\n", c.n, c.m, c.s)
		expR, expC := solve(c.n, c.m, c.s)
		exp := fmt.Sprintf("%d %d", expR, expC)
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
