package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(input string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scanner.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(scanner.Text())
	}
	n, err := nextInt()
	if err != nil {
		return "", err
	}
	m, err := nextInt()
	if err != nil {
		return "", err
	}
	short := make(map[[2]int]bool)
	maxX, maxY := 0, 0
	for i := 0; i < n; i++ {
		x, _ := nextInt()
		y, _ := nextInt()
		short[[2]int{x, y}] = true
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}
	query := make([][2]int, m)
	for i := 0; i < m; i++ {
		a, _ := nextInt()
		b, _ := nextInt()
		query[i] = [2]int{a, b}
		if a > maxX {
			maxX = a
		}
		if b > maxY {
			maxY = b
		}
	}
	maxX++
	maxY++
	win := make([][]bool, maxX+1)
	for i := range win {
		win[i] = make([]bool, maxY+1)
	}
	for x := 0; x <= maxX; x++ {
		for y := 0; y <= maxY; y++ {
			if x == 0 && y == 0 {
				win[x][y] = false
				continue
			}
			if short[[2]int{x, y}] {
				win[x][y] = false
				continue
			}
			res := false
			for i := 1; i <= x && !res; i++ {
				if !win[x-i][y] {
					res = true
				}
			}
			for j := 1; j <= y && !res; j++ {
				if !win[x][y-j] {
					res = true
				}
			}
			win[x][y] = res
		}
	}
	var out strings.Builder
	for i := 0; i < m; i++ {
		a, b := query[i][0], query[i][1]
		if a <= maxX && b <= maxY && win[a][b] {
			out.WriteString("WIN\n")
		} else {
			out.WriteString("LOSE\n")
		}
	}
	return strings.TrimSpace(out.String()), nil
}

type Test struct{ input string }

func genTests() []Test {
	rand.Seed(4)
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		s := make(map[[2]int]bool)
		for j := 0; j < n; j++ {
			x := rand.Intn(5)
			y := rand.Intn(5)
			s[[2]int{x, y}] = true
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		for j := 0; j < m; j++ {
			a := rand.Intn(5)
			b := rand.Intn(5)
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		tests = append(tests, Test{sb.String()})
	}
	// simple
	tests = append(tests, Test{"1 1\n0 0\n0 0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		exp, err := solve(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
