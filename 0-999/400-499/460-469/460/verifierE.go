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

func generatePoints(r int) []point {
	var v []point
	for i := -r; i <= r; i++ {
		for j := -r; j <= r; j++ {
			d2 := i*i + j*j
			if d2 <= r*r && d2 >= (r-1)*(r-1) {
				v = append(v, point{i, j})
			}
		}
	}
	sort.Slice(v, func(i, j int) bool {
		di := v[i].x*v[i].x + v[i].y*v[i].y
		dj := v[j].x*v[j].x + v[j].y*v[j].y
		return di > dj
	})
	return v
}

func dfs(v []point, n int, k int, cur []point, best *[]point, maxSum *int64) {
	if k == n {
		var sum int64
		for i := 0; i < len(cur); i++ {
			for j := i + 1; j < len(cur); j++ {
				dx := int64(cur[i].x - cur[j].x)
				dy := int64(cur[i].y - cur[j].y)
				sum += dx*dx + dy*dy
			}
		}
		if sum > *maxSum {
			*maxSum = sum
			tmp := make([]point, len(cur))
			copy(tmp, cur)
			*best = tmp
		}
		return
	}
	limit := len(v)
	cut := 30 - 2*n
	if cut < limit {
		limit = cut
	}
	for i := k; i < limit; i++ {
		cur = append(cur, v[i])
		dfs(v, n, i, cur, best, maxSum)
		cur = cur[:len(cur)-1]
	}
}

func expectedE(n, r int) (int64, []point) {
	v := generatePoints(r)
	var best []point
	var maxSum int64 = -1 << 60
	dfs(v, n, 0, nil, &best, &maxSum)
	return maxSum, best
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
	expectSum, expectPts := expectedE(n, r)
	if int64(nums[0]) != expectSum {
		return fmt.Errorf("expected sum %d got %d", expectSum, nums[0])
	}
	idx := 1
	for i := 0; i < n; i++ {
		x := nums[idx]
		y := nums[idx+1]
		if x != expectPts[i].x || y != expectPts[i].y {
			return fmt.Errorf("expected points %v got %v", expectPts, nums[1:])
		}
		idx += 2
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
