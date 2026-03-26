package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type point struct{ x, y int }

func crossProduct(o, a, b point) int {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

func convexHull(points []point) []point {
	sort.Slice(points, func(i, j int) bool {
		if points[i].x == points[j].x {
			return points[i].y < points[j].y
		}
		return points[i].x < points[j].x
	})
	var hull []point
	for _, p := range points {
		for len(hull) >= 2 && crossProduct(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	t := len(hull)
	for i := len(points) - 2; i >= 0; i-- {
		p := points[i]
		for len(hull) > t && crossProduct(hull[len(hull)-2], hull[len(hull)-1], p) <= 0 {
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, p)
	}
	if len(hull) > 1 {
		hull = hull[:len(hull)-1]
	}
	return hull
}

func expectedE(n, r int) int64 {
	var points []point
	for x := -r; x <= r; x++ {
		for y := -r; y <= r; y++ {
			if x*x+y*y <= r*r {
				points = append(points, point{x, y})
			}
		}
	}
	hull := convexHull(points)

	const off = 240
	const dim = 481
	// dp[sx][sy] for current layer; -1 means unreachable
	cur := make([][]int, dim)
	nxt := make([][]int, dim)
	for i := range cur {
		cur[i] = make([]int, dim)
		nxt[i] = make([]int, dim)
		for j := range cur[i] {
			cur[i][j] = -1
			nxt[i][j] = -1
		}
	}
	cur[off][off] = 0

	for step := 0; step < n; step++ {
		for j := range nxt {
			for k := range nxt[j] {
				nxt[j][k] = -1
			}
		}
		for sx := 0; sx < dim; sx++ {
			for sy := 0; sy < dim; sy++ {
				if cur[sx][sy] == -1 {
					continue
				}
				for _, p := range hull {
					nsx := sx + p.x
					nsy := sy + p.y
					if nsx < 0 || nsx >= dim || nsy < 0 || nsy >= dim {
						continue
					}
					nval := cur[sx][sy] + p.x*p.x + p.y*p.y
					if nval > nxt[nsx][nsy] {
						nxt[nsx][nsy] = nval
					}
				}
			}
		}
		cur, nxt = nxt, cur
	}

	var maxSum int64 = -1
	for sx := 0; sx < dim; sx++ {
		for sy := 0; sy < dim; sy++ {
			if cur[sx][sy] == -1 {
				continue
			}
			realSx := sx - off
			realSy := sy - off
			val := int64(n)*int64(cur[sx][sy]) - int64(realSx*realSx) - int64(realSy*realSy)
			if val > maxSum {
				maxSum = val
			}
		}
	}
	return maxSum
}

func runCase(bin string, n, r int) error {
	input := fmt.Sprintf("%d %d\n", n, r)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	var nums []int
	for scanner.Scan() {
		var x int
		fmt.Sscan(scanner.Text(), &x)
		nums = append(nums, x)
	}
	if len(nums) < 1+2*n {
		return fmt.Errorf("bad output")
	}
	expectSum := expectedE(n, r)
	if int64(nums[0]) != expectSum {
		return fmt.Errorf("expected sum %d got %d", expectSum, nums[0])
	}
	idx := 1
	for i := 0; i < n; i++ {
		x := nums[idx]
		y := nums[idx+1]
		if x*x+y*y > r*r {
			return fmt.Errorf("point (%d,%d) outside circle of radius %d", x, y, r)
		}
		idx += 2
	}
	// verify claimed sum matches actual pairwise distances
	var actualSum int64
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		pts[i] = point{nums[1+2*i], nums[2+2*i]}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := int64(pts[i].x - pts[j].x)
			dy := int64(pts[i].y - pts[j].y)
			actualSum += dx*dx + dy*dy
		}
	}
	if actualSum != int64(nums[0]) {
		return fmt.Errorf("claimed sum %d but actual pairwise sum is %d", nums[0], actualSum)
	}
	return nil
}

func generateCase(rng *rand.Rand) (int, int) {
	n := rng.Intn(7) + 2
	r := rng.Intn(10) + 1
	return n, r
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edges := []struct{ n, r int }{
		{2, 1}, {3, 2}, {5, 5},
	}
	for i, e := range edges {
		if err := runCase(bin, e.n, e.r); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		n, r := generateCase(rng)
		if err := runCase(bin, n, r); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
