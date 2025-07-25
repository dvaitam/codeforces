package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func countAtLeastK(arr []int, k int) int64 {
	left := 0
	sum := 0
	var res int64
	for right := 0; right < len(arr); right++ {
		sum += arr[right]
		for sum >= k {
			sum -= arr[left]
			left++
		}
		res += int64(left)
	}
	return res
}

func expectedE(r, c, k int, points [][2]int) int64 {
	grid := make([][]byte, r)
	for i := range grid {
		grid[i] = make([]byte, c)
	}
	for _, p := range points {
		grid[p[0]][p[1]] = 1
	}
	if r > c {
		ng := make([][]byte, c)
		for i := 0; i < c; i++ {
			ng[i] = make([]byte, r)
			for j := 0; j < r; j++ {
				ng[i][j] = grid[j][i]
			}
		}
		grid = ng
		r, c = c, r
	}
	col := make([]int, c)
	var ans int64
	for top := 0; top < r; top++ {
		for j := 0; j < c; j++ {
			col[j] = 0
		}
		for bottom := top; bottom < r; bottom++ {
			for j := 0; j < c; j++ {
				if grid[bottom][j] == 1 {
					col[j]++
				}
			}
			ans += countAtLeastK(col, k)
		}
	}
	return ans
}

func generateE(rng *rand.Rand) (string, int64) {
	r := rng.Intn(5) + 1
	c := rng.Intn(5) + 1
	n := rng.Intn(r*c) + 1
	k := rng.Intn(n) + 1
	used := map[[2]int]bool{}
	points := make([][2]int, 0, n)
	for len(points) < n {
		x := rng.Intn(r)
		y := rng.Intn(c)
		p := [2]int{x, y}
		if !used[p] {
			used[p] = true
			points = append(points, p)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", r, c, n, k)
	for _, p := range points {
		fmt.Fprintf(&sb, "%d %d\n", p[0]+1, p[1]+1)
	}
	exp := expectedE(r, c, k, points)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(46))
	for i := 0; i < 100; i++ {
		input, exp := generateE(rng)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", i+1, err, out.String())
			return
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(exp) {
			fmt.Printf("case %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed")
}
