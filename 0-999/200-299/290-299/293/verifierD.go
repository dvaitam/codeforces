package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (one line per test).
const embeddedTestcases = `5 -4 -5 -2 4 -5 4 3 0 -3 2
5 -4 4 2 -1 -2 1 1 3 0 -2
5 -4 -5 -3 0 5 4 4 2 -4 3
6 1 0 1 -5 -4 1 4 5 -3 -3 -5 -4
6 -2 -5 5 -1 -1 1 -4 0 -1 0 4 5
4 3 3 -5 0 3 3 -5 -3
6 -5 5 -1 0 -3 0 -4 4 3 1 -4 -1
3 3 -4 4 -3 2 -5
4 0 -4 3 -2 -3 5 -2 1
4 4 2 0 -3 -1 0 3 -2
4 2 0 -3 -1 3 5 -5 -4
5 -1 0 -5 -4 4 -5 -1 -4 -3 -1
4 3 3 0 2 -4 1 3 3
6 -4 -4 4 2 0 -3 -1 -2 5 1 5 -1
3 -5 -5 5 -1 -4 5
5 0 3 5 -1 -5 -2 3 -3 -1 4
3 1 1 -1 4 -2 -1
5 3 2 2 -4 1 0 -3 3 -5 -2
6 1 -5 -2 -4 -1 -4 -3 -3 2 -5 0 3
6 -4 3 -2 2 -1 -1 1 1 2 -5 -3 3
5 -2 -2 -3 -2 -1 -3 -3 3 -1 3
3 -3 -2 1 -3 -2 3
4 4 5 -3 -2 3 -1 -5 2
4 -5 1 -3 1 3 1 -2 1
6 -4 -3 -4 3 -5 -4 -5 -2 1 2 5 -1
6 0 5 1 -1 -4 -3 4 -2 -2 -5 -5 3
3 2 -2 2 5 -3 4
5 -5 4 -3 0 5 0 0 -1 3 5
5 -3 -5 2 -2 4 1 0 2 -1 2
4 2 -1 0 0 5 4 1 1
6 3 -1 5 4 -2 -2 1 -2 -1 3 -3 3
4 4 3 3 3 5 2 0 -1
6 -1 2 5 3 -5 -2 -1 2 -4 -5 -2 4
3 -1 -5 -2 -3 4 -3
5 4 -1 -5 5 5 -5 -5 4 -3 -1
4 1 2 0 0 -4 4 1 -4
3 -3 -5 4 3 -1 0
4 1 -2 0 -3 -3 2 5 -3
6 4 5 -2 -5 -5 1 4 -5 5 5 5 -3
3 1 -1 0 5 5 4
4 -5 0 -3 0 5 -1 2 2
5 -1 5 -4 2 -3 -4 3 -1 3 1
6 1 5 -3 4 1 -3 5 5 -4 4 -4 -3
6 -1 5 -4 5 1 -1 -1 -4 5 0 3 3
3 -5 -5 4 -5 -3 -4
6 1 2 -5 0 5 -3 -4 1 5 4 -2 -5
4 1 -5 -5 4 1 -3 3 -1
6 1 3 -5 2 3 -3 3 -3 5 -2 4 -3
3 4 2 -1 -1 -1 -2
6 5 -5 5 -5 3 -4 -5 5 -1 -3 -3 3
4 -4 0 -5 0 -3 0 4 -4
3 0 -2 -1 2 0 5
6 2 3 3 -4 -4 0 2 1 1 3 2 3
6 3 -4 -2 -3 5 1 3 0 -3 -4 -3 -3
3 1 5 4 -4 5 -4
5 4 -1 0 5 2 1 4 0 5 4
6 3 0 -2 3 -5 1 0 0 3 -5 -2 4
3 2 -4 0 1 4 -2
3 3 3 -4 -5 4 0
3 0 -3 0 -1 2 2
6 -5 5 -2 5 5 2 0 5 -3 -5 2 3
4 1 -1 5 4 2 -2 -5 2
4 4 2 0 4 0 1 4 5
4 3 1 1 -2 0 1 -1 -2
4 -1 -2 3 4 0 2 -2 4
5 4 5 -3 0 -2 -3 -1 3 -3 3
3 2 2 -2 -3 4 0
4 -4 3 2 -2 -1 1 -1 3
4 0 1 0 -5 -5 -2 -2 4
4 3 2 -4 -1 3 -1 2 1
5 5 1 -5 -4 -5 0 -3 -1 4 3
6 -1 3 1 2 0 -2 5 -4 -4 5 -4 3
3 5 4 3 0 4 -4
4 -5 -5 -4 -1 -4 4 1 -3
4 -1 -1 -3 -3 4 -2 -3 -3
4 -1 2 -1 -3 4 -1 -2 -4
4 5 1 -3 -2 4 -5 0 3
6 5 -3 -1 -4 4 -1 4 -2 -5 -3 -2 4
3 1 -3 -3 2 3 -5
5 -2 5 -3 2 0 5 2 -2 5 5
5 5 -2 -4 0 2 -1 0 3 4 -1
4 -1 3 1 -4 -3 1 -4 4
3 3 -3 2 -3 1 -1
4 -3 -4 -1 2 0 1 3 -5
5 1 1 3 -3 3 4 -3 5 -5 -1
6 3 4 1 -2 2 -2 1 -4 -4 3 1 -4
6 5 -3 -5 -2 -2 -2 -4 4 1 1 -4 3
4 3 -2 -2 2 -4 -2 -4 3
4 4 -2 3 4 3 5 4 1
6 5 -1 5 2 5 5 4 -5 0 0 -2 5
3 -2 2 -2 3 5 -2
3 -3 -2 2 2 -4 4
5 5 -1 1 0 5 -1 5 -4 4 -4
6 0 4 4 5 -4 -2 -3 1 -2 -5 2 -1
4 3 5 0 5 -3 -4 3 0
5 5 -4 0 -2 -2 -1 3 1 5 1
6 -5 -3 2 -4 -1 5 5 -3 3 -2 4 2
5 -3 4 -3 4 1 -1 -2 -4 -3 -3
6 3 1 -3 -5 -1 2 4 -4 5 0 5 -4
3 2 -2 3 -5 -5 5`

const ZERO = 1100000

var (
	fromArr = make([]int, 2*ZERO+1)
	toArr   = make([]int, 2*ZERO+1)
	x       []int
	y       []int
)

func det(a, b, c, d int64) int64 {
	return a*d - b*c
}

func ensureXY(n int) {
	if len(x) < n+1 {
		x = make([]int, n+1)
		y = make([]int, n+1)
	}
}

func calc(n int) float64 {
	var area int64
	x[n] = x[0]
	y[n] = y[0]
	for i := 0; i < n; i++ {
		area += det(int64(x[i]), int64(y[i]), int64(x[i+1]), int64(y[i+1]))
	}
	if area > 0 {
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			x[i], x[j] = x[j], x[i]
			y[i], y[j] = y[j], y[i]
		}
		x[n] = x[0]
		y[n] = y[0]
	}
	minX, maxX := 2*ZERO, 0
	for i := 0; i < n; i++ {
		if x[i] < minX {
			minX = x[i]
		}
		if x[i] > maxX {
			maxX = x[i]
		}
	}
	for tx := minX; tx <= maxX; tx++ {
		fromArr[tx] = 0
		toArr[tx] = 0
	}
	for i := 0; i < n; i++ {
		xi, yi := x[i], y[i]
		xj, yj := x[i+1], y[i+1]
		if xi < xj {
			dx := xj - xi
			for tx := xi; tx <= xj; tx++ {
				ty := (yj*(tx-xi) + yi*(xj-tx)) / dx
				toArr[tx] = ty
			}
		} else if xi > xj {
			dx := xi - xj
			for tx := xj; tx <= xi; tx++ {
				ty := (yi*(tx-xj) + yj*(xi-tx) + dx - 1) / dx
				fromArr[tx] = ty
			}
		} else {
			y0, y1 := yi, yj
			if y0 > y1 {
				y0, y1 = y1, y0
			}
			fromArr[xi] = y0
			toArr[xi] = y1
		}
	}
	var ans, sum, sumSq, cnt float64
	for tx := minX; tx <= maxX; tx++ {
		cur := toArr[tx] - fromArr[tx] + 1
		if cur <= 0 {
			continue
		}
		fcur := float64(cur)
		ftx := float64(tx)
		ans += fcur*cnt*ftx*ftx + fcur*sumSq - 2*fcur*ftx*sum
		sum += fcur * ftx
		sumSq += fcur * ftx * ftx
		cnt += fcur
	}
	return ans * 2 / (cnt * (cnt - 1))
}

func solve(points [][2]int) string {
	n := len(points)
	ensureXY(n)
	for i, p := range points {
		x[i] = p[0] + ZERO
		y[i] = p[1] + ZERO
	}
	val := calc(n)
	for i := 0; i < n; i++ {
		x[i], y[i] = y[i], x[i]
	}
	val += calc(n)
	return fmt.Sprintf("%.6f", val/2)
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: bad n\n", idx+1)
			os.Exit(1)
		}
		if len(fields) != 1+2*n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d coords got %d\n", idx+1, 2*n, len(fields)-1)
			os.Exit(1)
		}
		points := make([][2]int, n)
		for i := 0; i < n; i++ {
			xv, _ := strconv.Atoi(fields[1+2*i])
			yv, _ := strconv.Atoi(fields[2+2*i])
			points[i] = [2]int{xv, yv}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for _, p := range points {
			input.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		expected := solve(points)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
