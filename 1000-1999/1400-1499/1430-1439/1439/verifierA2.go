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

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func copyGrid(g [][]int) [][]int {
	res := make([][]int, len(g))
	for i := range g {
		res[i] = append([]int(nil), g[i]...)
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func check(bin string, n, m int, grid [][]int, limit int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte('\n')
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return fmt.Errorf("exec error: %v output=%s", err, out)
	}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(sc.Text())
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k < 0 || k > limit {
		return fmt.Errorf("k out of range")
	}
	for t := 0; t < k; t++ {
		xs := make([]int, 3)
		ys := make([]int, 3)
		for i := 0; i < 3; i++ {
			if !sc.Scan() {
				return fmt.Errorf("not enough numbers")
			}
			xi, err := strconv.Atoi(sc.Text())
			if err != nil {
				return fmt.Errorf("invalid int")
			}
			if !sc.Scan() {
				return fmt.Errorf("not enough numbers")
			}
			yi, err := strconv.Atoi(sc.Text())
			if err != nil {
				return fmt.Errorf("invalid int")
			}
			xi--
			yi--
			if xi < 0 || xi >= n || yi < 0 || yi >= m {
				return fmt.Errorf("coord out")
			}
			xs[i], ys[i] = xi, yi
		}
		// distinct
		if xs[0] == xs[1] && ys[0] == ys[1] || xs[0] == xs[2] && ys[0] == ys[2] || xs[1] == xs[2] && ys[1] == ys[2] {
			return fmt.Errorf("duplicate cells")
		}
		// same 2x2
		xmn, xmx := xs[0], xs[0]
		ymn, ymx := ys[0], ys[0]
		for i := 1; i < 3; i++ {
			xmn = min(xmn, xs[i])
			xmx = max(xmx, xs[i])
			ymn = min(ymn, ys[i])
			ymx = max(ymx, ys[i])
		}
		if xmx-xmn > 1 || ymx-ymn > 1 {
			return fmt.Errorf("not in one 2x2")
		}
		for i := 0; i < 3; i++ {
			grid[xs[i]][ys[i]] ^= 1
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 0 {
				return fmt.Errorf("not zero")
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/bin")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 0; t < 100; t++ {
		n := rand.Intn(3) + 2
		m := rand.Intn(3) + 2
		g := make([][]int, n)
		for i := 0; i < n; i++ {
			g[i] = make([]int, m)
			for j := 0; j < m; j++ {
				g[i][j] = rand.Intn(2)
			}
		}
		if err := check(bin, n, m, copyGrid(g), n*m); err != nil {
			fmt.Printf("test %d failed: %v\n", t, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
