package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Point structure copied from solution
type Point struct {
	x, y int64
}

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isGood(dx, dy int64, pts []Point, sumX, sumY int64, n int) bool {
	m := make(map[int64]int, n)
	target := sumX*dx + sumY*dy
	for _, p := range pts {
		val := (p.x*dx+p.y*dy)*int64(n) - target
		m[val]++
	}
	for k, v := range m {
		if m[-k] != v {
			return false
		}
	}
	return true
}

func expectedF(pts []Point) int {
	n := len(pts)
	var sumX, sumY int64
	for _, p := range pts {
		sumX += p.x
		sumY += p.y
	}
	sorted := make([]Point, n)
	copy(sorted, pts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].x == sorted[j].x {
			return sorted[i].y < sorted[j].y
		}
		return sorted[i].x < sorted[j].x
	})
	candX := sorted[0].x + sorted[n-1].x
	candY := sorted[0].y + sorted[n-1].y
	symmetric := true
	for i := 0; i < n; i++ {
		if sorted[i].x+sorted[n-1-i].x != candX || sorted[i].y+sorted[n-1-i].y != candY {
			symmetric = false
			break
		}
	}
	if symmetric {
		return -1
	}
	type pair struct{ x, y int64 }
	dirs := make(map[pair]struct{})
	idx := []int{0, 1}
	if n == 1 {
		return 0
	}
	for _, i := range idx {
		if i >= n {
			continue
		}
		for j := i + 1; j < n; j++ {
			ux := (pts[i].x+pts[j].x)*int64(n) - 2*sumX
			uy := (pts[i].y+pts[j].y)*int64(n) - 2*sumY
			if ux == 0 && uy == 0 {
				continue
			}
			dx := -uy
			dy := ux
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			dirs[pair{dx, dy}] = struct{}{}
		}
	}
	count := 0
	for d := range dirs {
		if isGood(d.x, d.y, pts, sumX, sumY, n) {
			count++
		}
	}
	return count
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	pts := make([]Point, 0, n)
	used := map[[2]int]struct{}{}
	for len(pts) < n {
		x := rng.Intn(7) - 3
		y := rng.Intn(7) - 3
		key := [2]int{x, y}
		if _, ok := used[key]; ok {
			continue
		}
		used[key] = struct{}{}
		pts = append(pts, Point{int64(x), int64(y)})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, p := range pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return testCase{input: sb.String(), expected: expectedF(pts)}
}

type testCase struct {
	input    string
	expected int
}

func run(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %d\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
